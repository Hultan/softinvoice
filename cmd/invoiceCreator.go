package main

import (
	"fmt"
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/hultan/softinvoice/internal/database"
	"github.com/hultan/softteam/framework"
	"github.com/jung-kurt/gofpdf"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"io/ioutil"
	"os"
)

type InvoiceCreator struct {
	Invoice *database.Invoice
}

func NewInvoiceCreator(invoice *database.Invoice) *InvoiceCreator {
	creator := new(InvoiceCreator)
	creator.Invoice = invoice
	return creator
}

func (i *InvoiceCreator) CreatePDF(path string) {
	_, imagePath := i.CreatePNG()

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	//pdf.SetFont("Arial", "B", 16)
	//pdf.Cell(40, 10, "Hello world!!!")
	pdf.Image(imagePath,15,30,180,247,false,"",0,"")

	err := pdf.OutputFileAndClose(path)
	errorCheck(err)

	// Clean up image
	err = os.Remove(imagePath)
	errorCheck(err)
}

func (i *InvoiceCreator) CreatePNG() (*gdk.Pixbuf, string) {
	// Load the image
	fw := framework.NewFramework()
	filePath := fw.Resource.GetResourcePath("empty_invoice.png")
	image, err :=gdk.PixbufNewFromFile(filePath)
	if err!=nil {
		panic(err)
	}

	// Get the image format
	var format cairo.Format
	if image.GetHasAlpha () {
		format = cairo.FORMAT_ARGB32
	} else {
		format = cairo.FORMAT_RGB24
	}

	// Create a surface to draw on
	width := image.GetWidth ()
	height := image.GetHeight()
	surface := cairo.CreateImageSurface (format, width, height)
	if surface==nil {
		panic("surface is nil")
	}
	cr := cairo.Create(surface)

	// Load the image into the surface
	gtk.GdkCairoSetSourcePixBuf(cr,image,0,0)
	cr.Paint()

	// Fill in invoice text
	i.FillInvoiceTextPNG(cr)

	// Save the image
	file, err := ioutil.TempFile("/tmp","se_softteam_invoice_*.png")
	errorCheck(err)
	err = surface.WriteToPNG(file.Name())
	errorCheck(err)

	returnImage, err :=gdk.PixbufNewFromFile(file.Name())

	// Clean up
	surface = nil
	cr = nil

	return returnImage, file.Name()
}

func (i *InvoiceCreator) FillInvoiceTextPNG(cr *cairo.Context) {
	// Header : Left
	i.WriteTextOnPNG(cr, 15, cairo.FONT_WEIGHT_NORMAL, 380,216,i.Invoice.Date.Format(constDateLayout), true)
	i.WriteTextOnPNG(cr, 15, cairo.FONT_WEIGHT_NORMAL, 380,236,fmt.Sprintf("%d",i.Invoice.Number), true)
	i.WriteTextOnPNG(cr, 15, cairo.FONT_WEIGHT_NORMAL, 380,257,i.Invoice.CustomerNumber, true)
	i.WriteTextOnPNG(cr, 15, cairo.FONT_WEIGHT_NORMAL, 380,287,i.Invoice.CustomerReference, true)
	i.WriteTextOnPNG(cr, 15, cairo.FONT_WEIGHT_NORMAL, 380,307,"Per Hultqvist", true)

	// Header : Right
	i.WriteTextOnPNG(cr, 15, cairo.FONT_WEIGHT_NORMAL, 450,236,i.Invoice.CustomerName,false)
	i.WriteTextOnPNG(cr, 15, cairo.FONT_WEIGHT_NORMAL, 450,256,i.Invoice.CustomerAddress,false)
	i.WriteTextOnPNG(cr, 15, cairo.FONT_WEIGHT_NORMAL, 450,276,i.Invoice.CustomerPostalAddress,false)

	// Rows
	p := message.NewPrinter(language.Swedish)
	var sumExclVAT float32 = 0.0

	for index, invoiceRow := range i.Invoice.Rows {
		offset := float64(index) * 50
		i.WriteTextOnPNG(cr, 15, cairo.FONT_WEIGHT_NORMAL, 110,400 + offset,invoiceRow.Text,false)
		i.WriteTextOnPNG(cr, 15, cairo.FONT_WEIGHT_NORMAL, 110,420 + offset,invoiceRow.Name,false)
		i.WriteTextOnPNG(cr, 15, cairo.FONT_WEIGHT_NORMAL, 523,420 + offset,p.Sprintf("%.1f",invoiceRow.Amount), true)
		i.WriteTextOnPNG(cr, 15, cairo.FONT_WEIGHT_NORMAL, 600,420 + offset,p.Sprintf("%.0f",invoiceRow.Price), true)
		i.WriteTextOnPNG(cr, 15, cairo.FONT_WEIGHT_NORMAL, 710,420 + offset,p.Sprintf("%.2f",invoiceRow.Total), true)

		sumExclVAT += invoiceRow.Total
	}

	vat:=sumExclVAT*0.25
	sumInclVAT:=sumExclVAT*1.25
	toPay := float32(int(sumInclVAT))
	rounded:=float32(int(sumInclVAT)) - sumInclVAT
	i.WriteTextOnPNG(cr, 15, cairo.FONT_WEIGHT_NORMAL, 710,789 ,p.Sprintf("%.2f",sumExclVAT), true)
	i.WriteTextOnPNG(cr, 15, cairo.FONT_WEIGHT_NORMAL, 710,807 ,p.Sprintf("%.2f",vat), true)
	i.WriteTextOnPNG(cr, 15, cairo.FONT_WEIGHT_NORMAL, 710,832 ,p.Sprintf("%.2f",sumInclVAT), true)
	i.WriteTextOnPNG(cr, 15, cairo.FONT_WEIGHT_NORMAL, 710,850 ,p.Sprintf("%.2f",rounded), true)

	i.WriteTextOnPNG(cr, 16, cairo.FONT_WEIGHT_NORMAL, 190,883 ,"5689-1849", false)
	i.WriteTextOnPNG(cr, 16, cairo.FONT_WEIGHT_NORMAL, 405,883 ,i.Invoice.DueDate.Format(constDateLayout), false)
	i.WriteTextOnPNG(cr, 16, cairo.FONT_WEIGHT_BOLD, 710,883 ,p.Sprintf("%.2f",toPay), true)

}

func (i *InvoiceCreator) WriteTextOnPNG(cr *cairo.Context, fontSize float64, fontWeight cairo.FontWeight, x float64,y float64, text string, rightJustify bool ) {
	// Write text
	cr.SetSourceRGB(0,0,0)
	cr.SetFontSize(fontSize)
	cr.SelectFontFace("Cantarell",cairo.FONT_SLANT_NORMAL, fontWeight)
	if rightJustify {
		te := cr.TextExtents(text)
		cr.MoveTo(x - te.Width,y)
		cr.ShowText(text)
	} else {
		cr.MoveTo(x,y)
		cr.ShowText(text)
	}
}

func (i *InvoiceCreator) GetTextSize(cr *cairo.Context, text string) cairo.TextExtents {
	te := cr.TextExtents(text)
	return te
}