package main

import "C"
import (
	"fmt"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/hultan/softinvoice/internal/database"
	"github.com/hultan/softteam/framework"
	"os"
	"strconv"
)

type MainForm struct {
	window    *gtk.ApplicationWindow
	popupMenu *PopupMenu
	treeView  *gtk.TreeView
	invoices  []database.Invoice
}

func NewMainForm() *MainForm {
	mainForm := new(MainForm)
	return mainForm
}

func (m *MainForm) OpenMainForm(app *gtk.Application) {
	// Initialize gtk
	gtk.Init(&os.Args)

	// Create a new gtk helper
	fw := framework.NewFramework()

	builder, err := fw.Gtk.CreateBuilder("main.glade")
	if err!=nil {
		panic("glade file not found!")
	}
	softInvoice.builder = builder

	// Get the main window from the glade file
	window := softInvoice.builder.GetObject("main_window").(*gtk.ApplicationWindow)
	m.window = window

	// Set up main window
	window.SetApplication(app)
	title := fmt.Sprintf("%s - %s - [Database : %s]", applicationName, applicationVersion, softInvoice.database.GetDatabaseName())
	window.SetTitle(title)
	window.SetDefaultSize(1024, 768)

	// Hook up the destroy event
	_ = window.Connect("destroy", func() {
		m.CloseMainForm()
	})

	// Get the new invoice button
	button := softInvoice.builder.GetObject("newinvoice_button").(*gtk.ToolButton)

	// Hook up the clicked event for the new invoice button
	_ = button.Connect("clicked", func() {
		softInvoice.invoiceForm.OpenInvoiceForm(nil)
	})

	err = m.LoadInvoiceList()
	errorCheck(err)

	m.popupMenu = NewPopupMenu(softInvoice, m)

	// Show the main window
	window.ShowAll()
}

func (m *MainForm) CloseMainForm() {
	softInvoice.CleanUp()
}

func (m *MainForm) LoadInvoiceList() error {
	// Get all invoices from the database
	invoices, err := softInvoice.database.GetAllInvoices()
	if err != nil {
		return err
	}
	m.invoices = invoices

	// Get the treeview from the builder
	treeView := softInvoice.builder.GetObject("invoice_treeview").(*gtk.TreeView)
	m.treeView = treeView

	// Create a new list store
	listStore, err := gtk.ListStoreNew(
		glib.TYPE_STRING, // Nummer
		glib.TYPE_STRING, // Datum
		glib.TYPE_STRING, // Förfalodatum
		glib.TYPE_STRING, // Kund
		glib.TYPE_STRING, // Belopp
		glib.TYPE_STRING) // Background color (credit)
	if err != nil {
		return err
	}

	// Fill list store
	for _, invoice := range invoices {
		iter := listStore.Append()
		_ = listStore.Set(iter, []int{0, 1, 2, 3, 4, 5}, []interface{}{
			fmt.Sprintf("%d", invoice.Number),
			invoice.Date.Format(constDateLayout),
			invoice.DueDate.Format(constDateLayout),
			invoice.CustomerName,
			fmt.Sprintf("%.2f", invoice.Amount),
			m.GetInvoiceColor(&invoice)})
	}

	// Set model and hook up row activated signal
	treeView.SetModel(listStore)
	_ = treeView.Connect("row_activated", m.OnInvoiceClicked)

	return nil
}

//
// Signal handlers
//

func (m *MainForm) OnInvoiceClicked(treeView *gtk.TreeView, softInvoice *SoftInvoice) {
	invoice := m.GetSelectedInvoice(treeView)
	if invoice == nil {
		return
	}

	softInvoice.previewWForm.OpenPreviewForm(softInvoice, invoice)
	//creator := NewInvoiceCreator(invoice)
	//creator.CreatePDF("/home/per/temp/test.pdf")
}

//
// Misc functions
//

func (m *MainForm) GetSelectedInvoice(treeView *gtk.TreeView) *database.Invoice {
	selection, err := treeView.GetSelection()
	if err != nil {
		return nil
	}
	model, iter, ok := selection.GetSelected()
	if ok {
		value, err := model.(*gtk.TreeModel).GetValue(iter, listStoreColumnInvoiceNumber)
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

func (m *MainForm) GetInvoiceColor(invoice *database.Invoice) string {
	if invoice.Credit {
		return "ORANGE"
	} else {
		return "WHITE"
	}
}
