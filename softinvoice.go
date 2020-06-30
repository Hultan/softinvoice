package main

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/hultan/softteam-invoice/database"
	gtkhelper "github.com/hultan/softteam/gtk"
)

// The SoftInvoice application type
type SoftInvoice struct {
	application    *gtk.Application
	mainForm       *MainForm
	invoiceForm    *InvoiceForm
	invoiceRowForm *InvoiceRowForm
	previewWForm   *PreviewForm
	helper         *gtkhelper.GtkHelper

	database       *database.Database
}

// Create a new SoftInvoice object
func NewSoftInvoice(app *gtk.Application) *SoftInvoice {
	softInvoice := new(SoftInvoice)
	softInvoice.invoiceForm = NewInvoiceForm()
	softInvoice.invoiceRowForm = NewInvoiceRowForm()
	softInvoice.previewWForm = NewPreviewForm()
	softInvoice.database = new(database.Database)
	softInvoice.application = app
	return softInvoice
}

func (s *SoftInvoice) CleanUp() {
	// Destroy the invoice window if it has been created
	if s.invoiceRowForm != nil && s.invoiceRowForm.window != nil {
		s.invoiceRowForm.window.Destroy()
	}

	// Destroy the invoice window if it has been created
	if s.invoiceForm != nil && s.invoiceForm.window != nil {
		s.invoiceForm.window.Destroy()
	}

	// Destroy the preview window if it has been created
	if s.previewWForm != nil && s.previewWForm.window != nil {
		s.previewWForm.window.Destroy()
	}

	// Close the database
	s.database.CloseDatabase()
}