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

type MainWindow struct {
	window    *gtk.ApplicationWindow
	treeView  *gtk.TreeView
	popupMenu *PopupMenu
	invoices  []database.Invoice
}

func NewMainWindow() *MainWindow {
	mainWindow := new(MainWindow)
	return mainWindow
}

func (m *MainWindow) OpenMainWindow(app *gtk.Application, softInvoice *SoftInvoice) {
	// Initialize gtk
	gtk.Init(&os.Args)

	// Create a new gtk helper
	softInvoice.helper = gtkhelper.GtkHelperNewFromFile("resources/main.glade")
	// Get the main window from the glade file
	mainWindow, err := softInvoice.helper.GetApplicationWindow("main_window")
	errorCheck(err)

	m.window = mainWindow

	// Set up main window
	mainWindow.SetApplication(app)
	title := fmt.Sprintf("SoftInvoice - [Database : %s]", database.DatabaseName)
	mainWindow.SetTitle(title)
	mainWindow.SetDefaultSize(800, 600)

	// Hook up the destroy event
	mainWindow.Connect("destroy", func() {
		m.CloseMainWindow(softInvoice)
	})

	// Get the new invoice button
	button, err := softInvoice.helper.GetToolButton("newinvoice_button")
	errorCheck(err)

	// Hook up the clicked event for the new invoice button
	button.Connect("clicked", func() {
		softInvoice.invoiceWindow.OpenInvoiceWindow(softInvoice)
	})

	err = m.loadInvoiceList(softInvoice)
	if err != nil {
		// Failed to load invoice list
	}

	m.popupMenu = NewPopupMenu(softInvoice, m)

	// Show the main window
	mainWindow.ShowAll()
}

func (m *MainWindow) CloseMainWindow(softInvoice *SoftInvoice) {
	// Destroy the invoice window if it has been created
	if softInvoice.invoiceWindow != nil && softInvoice.invoiceWindow.window != nil {
		softInvoice.invoiceWindow.window.Destroy()
	}

	// Destroy the preview window if it has been created
	if softInvoice.previewWindow != nil && softInvoice.previewWindow.window != nil {
		softInvoice.previewWindow.window.Destroy()
	}

	// Close the database
	softInvoice.database.CloseDatabase()
}

func (m *MainWindow) loadInvoiceList(softInvoice *SoftInvoice) error {
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

func (m *MainWindow) invoiceClicked(treeView *gtk.TreeView, path *gtk.TreePath, column *gtk.TreeViewColumn, softInvoice *SoftInvoice) {
	invoice := m.getSelectedInvoice(treeView)
	if invoice == nil {
		return
	}

	softInvoice.previewWindow.OpenPreviewWindow(softInvoice, invoice)
	//creator := NewInvoiceCreator(invoice)
	//creator.CreatePDF("/home/per/temp/test.pdf")
}

func (m *MainWindow) getSelectedInvoice(treeView *gtk.TreeView) *database.Invoice {
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

func (m *MainWindow) getColor(invoice *database.Invoice) string {
	if invoice.Credit {
		return "RED"
	} else {
		return "WHITE"
	}
}
