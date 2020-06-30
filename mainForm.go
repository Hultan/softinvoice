package main

import "C"
import (
	"fmt"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/hultan/softteam-invoice/database"
	gtkhelper "github.com/hultan/softteam/gtk"
	"os"
	"strconv"
)

type MainForm struct {
	window    *gtk.ApplicationWindow
	treeView  *gtk.TreeView
	popupMenu *PopupMenu
	invoices  []database.Invoice
}

func NewMainForm() *MainForm {
	mainForm := new(MainForm)
	return mainForm
}

func (m *MainForm) OpenMainForm(app *gtk.Application, softInvoice *SoftInvoice) {
	// Initialize gtk
	gtk.Init(&os.Args)

	// Create a new gtk helper
	softInvoice.helper = gtkhelper.GtkHelperNewFromFile("resources/main.glade")
	// Get the main window from the glade file
	window, err := softInvoice.helper.GetApplicationWindow("main_window")
	errorCheck(err)

	m.window = window

	// Set up main window
	window.SetApplication(app)
	title := fmt.Sprintf("SoftInvoice - [Database : %s]", database.DatabaseName)
	window.SetTitle(title)
	window.SetDefaultSize(800, 600)

	// Hook up the destroy event
	window.Connect("destroy", func() {
		m.CloseMainWindow(softInvoice)
	})

	// Get the new invoice button
	button, err := softInvoice.helper.GetToolButton("newinvoice_button")
	errorCheck(err)

	// Hook up the clicked event for the new invoice button
	button.Connect("clicked", func() {
		softInvoice.invoiceForm.OpenInvoiceForm(softInvoice)
	})

	err = m.loadInvoiceList(softInvoice)
	if err != nil {
		// Failed to load invoice list
	}

	m.popupMenu = NewPopupMenu(softInvoice, m)

	// Show the main window
	window.ShowAll()
}

func (m *MainForm) CloseMainWindow(softInvoice *SoftInvoice) {
	// Destroy the invoice window if it has been created
	if softInvoice.invoiceRowForm != nil && softInvoice.invoiceRowForm.window != nil {
		softInvoice.invoiceRowForm.window.Destroy()
	}

	// Destroy the invoice window if it has been created
	if softInvoice.invoiceForm != nil && softInvoice.invoiceForm.window != nil {
		softInvoice.invoiceForm.window.Destroy()
	}

	// Destroy the preview window if it has been created
	if softInvoice.previewWForm != nil && softInvoice.previewWForm.window != nil {
		softInvoice.previewWForm.window.Destroy()
	}

	// Close the database
	softInvoice.database.CloseDatabase()
}

func (m *MainForm) loadInvoiceList(softInvoice *SoftInvoice) error {
	invoices, err := softInvoice.database.GetAllInvoices()
	if err != nil {
		return err
	}
	m.invoices = invoices
	treeView, err := softInvoice.helper.GetTreeView("invoice_treeview")
	if err != nil {
		// Failed to get invoice treeview
		return err
	}
	m.treeView = treeView

	listStore, err := gtk.ListStoreNew(
		glib.TYPE_STRING, // Nummer
		glib.TYPE_STRING, // Datum
		glib.TYPE_STRING, // Förfalodatum
		glib.TYPE_STRING, // Kund
		glib.TYPE_STRING, // Belopp
		glib.TYPE_STRING) // Background color (credit)

	if err != nil {
		// Failed to create list store
		return err
	}

	for _, invoice := range invoices {
		iter := listStore.Append()
		//color:=m.getColor(&invoice)
		listStore.Set(iter, []int{0, 1, 2, 3, 4, 5}, []interface{}{
			fmt.Sprintf("%d", invoice.Number),
			invoice.Date.Format(constDateLayout),
			invoice.DueDate.Format(constDateLayout),
			invoice.CustomerName,
			fmt.Sprintf("%.2f", invoice.Amount),
			m.getColor(&invoice)})
	}

	treeView.SetModel(listStore)

	treeView.Connect("row_activated", m.invoiceClicked, softInvoice)

	return nil
}

func (m *MainForm) invoiceClicked(treeView *gtk.TreeView, path *gtk.TreePath, column *gtk.TreeViewColumn, softInvoice *SoftInvoice) {
	invoice := m.getSelectedInvoice(treeView)
	if invoice == nil {
		return
	}

	softInvoice.previewWForm.OpenPreviewForm(softInvoice, invoice)
	//creator := NewInvoiceCreator(invoice)
	//creator.CreatePDF("/home/per/temp/test.pdf")
}

func (m *MainForm) getSelectedInvoice(treeView *gtk.TreeView) *database.Invoice {
	selection, err := treeView.GetSelection()
	if err != nil {
		return nil
	}
	model, iter, ok := selection.GetSelected()
	if ok {
		value, err := model.(*gtk.TreeModel).GetValue(iter, liststoreColumnInvoiceNumber)
		if err != nil {
			return nil
		}
		invoiceNumberString, err := value.GetString()
		if err != nil {
			return nil
		}
		invoiceNumber, err := strconv.Atoi(invoiceNumberString)
		if err != nil {
			return nil
		}
		for _, invoice := range m.invoices {
			if invoice.Number == invoiceNumber {
				return &invoice
			}
		}
		return nil
	}

	return nil
}

func (m *MainForm) getColor(invoice *database.Invoice) string {
	if invoice.Credit {
		return "RED"
	} else {
		return "WHITE"
	}
}
