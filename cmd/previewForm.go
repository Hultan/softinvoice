package main

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/hultan/softteam-invoice/internal/database"
	"os"
)

type PreviewForm struct {
	window *gtk.Window
}

func NewPreviewForm() *PreviewForm {
	previewForm := new(PreviewForm)
	return previewForm
}

func (p *PreviewForm) OpenPreviewForm(softInvoice *SoftInvoice, invoice *database.Invoice) {
	// Check if it is the first time we open the preview window
	if softInvoice.previewWForm.window==nil {
		// Get the preview window from glade
		window, err := softInvoice.helper.GetWindow("preview_window")
		errorCheck(err)

		// Save a pointer to the preview window
		softInvoice.previewWForm.window = window

		// Set up the preview window
		window.SetApplication(softInvoice.application)
		window.HideOnDelete()
		window.SetModal(true)
		window.SetKeepAbove(true)
		window.SetPosition(gtk.WIN_POS_CENTER_ALWAYS)

		// Hook up the hide event
		_, err = window.Connect("hide", func() {
			p.ClosePreviewWindow()
		})
		errorCheck(err)

		// Show the window
		window.ShowAll()
	} else {
		// Show the window
		softInvoice.previewWForm.window.ShowAll()
	}

	image, err := softInvoice.helper.GetImage("invoice_preview")
	if err!=nil {
		panic(err)
	}
	creator := NewInvoiceCreator(invoice)
	pixBuf, path := creator.CreatePNG()
	image.SetFromPixbuf(pixBuf)

	// Clean up image
	err = os.Remove(path)
	errorCheck(err)
}

func (p *PreviewForm) ClosePreviewWindow() {

}