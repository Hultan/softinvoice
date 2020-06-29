package main

import (
	"fmt"
	"github.com/gotk3/gotk3/gtk"
)

type InvoiceWindow struct {
	window *gtk.Window
}

var isSaving = false

func NewInvoiceWindow() *InvoiceWindow {
	invoiceWindow := new(InvoiceWindow)
	return invoiceWindow
}

func (i *InvoiceWindow) OpenInvoiceWindow(softInvoice *SoftInvoice) {
	// Check if it is the first time we open the invoice window
	if softInvoice.invoiceWindow.window==nil {
		// Get the invoice window from glade
		window, err := softInvoice.helper.GetWindow("invoice_window")
		errorCheck(err)

		// Save a pointer to the invoice window
		softInvoice.invoiceWindow.window = window

		// Set up the invoice window
		window.SetApplication(softInvoice.application)
		window.HideOnDelete()
		window.SetModal(true)
		window.SetKeepAbove(true)
		window.SetPosition(gtk.WIN_POS_CENTER_ALWAYS)

		// Hook up the hide event
		window.Connect("hide", func() {
			i.CloseInvoiceWindow()
		})

		// Get the cancel button
		cancelButton, err := softInvoice.helper.GetButton("cancel_button")
		errorCheck(err)

		// Hook up the clicked event for the cancel button
		cancelButton.Connect("clicked", func() {
			window.Hide()
		})

		// Get the save button
		saveButton, err := softInvoice.helper.GetButton("save_button")
		errorCheck(err)

		// Hook up the clicked event for the cancel button
		saveButton.Connect("clicked", func() {
			isSaving = true
			window.Hide()
		})

		// Show the window
		window.ShowAll()
	} else {
		// Show the window
		softInvoice.invoiceWindow.window.ShowAll()
	}
}

func (i *InvoiceWindow) CloseInvoiceWindow() {
	if isSaving {
		isSaving = false
		i.saveInvoice()
	}
}

func (i *InvoiceWindow) saveInvoice() {
	fmt.Println("SAVE!")
}