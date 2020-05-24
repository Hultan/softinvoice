package main

type InvoiceWindow struct {
}

func NewInvoiceWindow() *InvoiceWindow {
	invoiceWindow := new(InvoiceWindow)
	return invoiceWindow
}

func (i *InvoiceWindow) OpenInvoiceWindow(softInvoice *SoftInvoice) {
	// Check if it is the first time we open the invoice window
	if softInvoice.invoiceWindow==nil {
		// Get the invoice window from glade
		window, err := softInvoice.helper.GetWindow("invoice_window")
		errorCheck(err)

		// Save a pointer to the invoice window
		softInvoice.invoiceWindow = window

		// Set up the invoice window
		window.SetApplication(softInvoice.application)
		window.HideOnDelete()
		window.SetModal(true)
		window.SetKeepAbove(true)

		// Hook up the hide event
		window.Connect("hide", func() {
			i.CloseInvoiceWindow()
		})

		// Get the cancel button
		button, err := softInvoice.helper.GetButton("cancel_button")
		errorCheck(err)

		// Hook up the clicked event for the cancel button
		button.Connect("clicked", func() {
			window.Hide()
		})

		// Show the window
		window.ShowAll()
	} else {
		// Show the window
		softInvoice.invoiceWindow.ShowAll()
	}
}

func (i *InvoiceWindow) CloseInvoiceWindow() {

}