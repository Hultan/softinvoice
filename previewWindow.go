package main

type PreviewWindow struct {

}

func NewPreviewWindow() *PreviewWindow {
	window := new(PreviewWindow)
	return window
}

func (p *PreviewWindow) OpenPreviewWindow(softInvoice *SoftInvoice, invoice *Invoice) {
	// Check if it is the first time we open the preview window
	if softInvoice.previewWindow==nil {
		// Get the preview window from glade
		window, err := softInvoice.helper.GetWindow("preview_window")
		errorCheck(err)

		// Save a pointer to the preview window
		softInvoice.previewWindow = window

		// Set up the preview window
		window.SetApplication(softInvoice.application)
		window.HideOnDelete()
		window.SetModal(true)
		window.SetKeepAbove(true)

		// Hook up the hide event
		window.Connect("hide", func() {
			p.ClosePreviewWindow()
		})

		// Show the window
		window.ShowAll()
	} else {
		// Show the window
		softInvoice.previewWindow.ShowAll()
	}

	image, err := softInvoice.helper.GetImage("invoice_preview")
	if err!=nil {
		panic(err)
	}
	creator := NewInvoiceCreator(invoice)
	pixbuf := creator.CreatePNG()
	image.SetFromPixbuf(pixbuf)
}

func (p *PreviewWindow) ClosePreviewWindow() {

}