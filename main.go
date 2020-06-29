package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	database "github.com/hultan/softteam-invoice/database"
	gtkhelper "github.com/hultan/softteam/gtk"
	"os"
)

const ApplicationId string = "se.softteam.invoice"
const ApplicationFlags glib.ApplicationFlags = glib.APPLICATION_FLAGS_NONE

// The SoftInvoice application type
type SoftInvoice struct {
	application   *gtk.Application
	invoiceWindow	*InvoiceWindow
	previewWindow	*PreviewWindow
	//invoiceWindow *gtk.Window
	//previewWindow *gtk.Window
	helper        *gtkhelper.GtkHelper
	database      *database.Database
}

func main() {
	// Create a new application
	app, err := gtk.ApplicationNew(ApplicationId, ApplicationFlags)
	errorCheck(err)

	// Create the SoftInvoice application object
	softInvoice := NewSoftInvoice(app)
	mainWindow := NewMainWindow()

	// Hook up the activate event handler
	app.Connect("activate", mainWindow.OpenMainWindow, softInvoice)

	// Start the application
	status := app.Run(os.Args)

	// Exit
	os.Exit(status)
}

// Create a new SoftInvoice object
func NewSoftInvoice(app *gtk.Application) *SoftInvoice {
	softInvoice := new(SoftInvoice)
	softInvoice.invoiceWindow = NewInvoiceWindow()
	softInvoice.previewWindow = NewPreviewWindow()
	softInvoice.database = new(database.Database)
	softInvoice.application = app
	return softInvoice
}

func errorCheck(err error) {
	if err != nil {
		panic(err)
	}
}

func softErrorCheck(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

//
//
//invoices, err := softInvoice.database.GetAllInvoices()
//softErrorCheck(err)
//fmt.Println(invoices[10].CustomerName)
//fmt.Println(invoices[10].rows[0].Price)
//
//invoice, err := softInvoice.database.GetInvoiceByNumber("1245")
//softErrorCheck(err)
//fmt.Println(invoice.CustomerName)
//fmt.Println(invoice.rows[0].Price)
//
//customers, err := softInvoice.database.GetAllCustomers()
//softErrorCheck(err)
//fmt.Println(customers[1].FancyName)
//
//customer, err := softInvoice.database.GetCustomerByNumber("1021")
//softErrorCheck(err)
//fmt.Println(customer.FancyName)
//
//// Should fail
//_, err = softInvoice.database.GetCustomerByNumber("1234")
//softErrorCheck(err)
//
//
//products, err := softInvoice.database.GetAllProducts()
//softErrorCheck(err)
//fmt.Println(len(products))
//
//product, err := softInvoice.database.GetProductByNumber("NOVA")
//softErrorCheck(err)
//fmt.Println(product.Name)
//
//// Should fail
//product, err = softInvoice.database.GetProductByNumber("FAIL")
//softErrorCheck(err)
//if product != nil {
//panic("This should have failed!")
//}
//
//var invoice2 = Invoice{}
//invoice2.Amount=100
//fmt.Println(invoice2.Credit)
//fmt.Println(invoice2.CreditInvoiceNumber)
//
//product2, err := softInvoice.database.GetProductByNumber("TEST")
//if err!=nil {
//panic(err)
//}
//product2.Price = 900
//err = softInvoice.database.UpdateProduct(product2)
//if err!=nil {
//panic(err)
//}
//product3, err := softInvoice.database.GetProductByNumber("TEST")
//if err!=nil {
//panic(err)
//}
//fmt.Println("Nytt pris : ", product3.Price)
