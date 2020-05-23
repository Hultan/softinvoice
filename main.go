package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	gtkhelper "github.com/hultan/softteam/gtk"
	"os"
)

const ApplicationId string = "se.softteam.invoice"
const ApplicationFlags glib.ApplicationFlags = glib.APPLICATION_FLAGS_NONE

// The SoftInvoice application type
type SoftInvoice struct {
	application   *gtk.Application
	invoiceWindow *gtk.Window
	helper        *gtkhelper.GtkHelper
	database      *database
}

func main() {
	// Create a new application
	app, err := gtk.ApplicationNew(ApplicationId, ApplicationFlags)
	errorCheck(err)

	// Create the SoftInvoice application object
	softInvoice := NewSoftInvoice(app)

	// Hook up the activate event handler
	app.Connect("activate", activate, softInvoice)

	// Start the application
	status := app.Run(os.Args)

	// Exit
	os.Exit(status)
}

func activate(app *gtk.Application, softInvoice *SoftInvoice) {
	// Initialize gtk
	gtk.Init(&os.Args)

	// Create a new gtk helper
	softInvoice.helper = gtkhelper.GtkHelperNewFromFile("resources/main.glade")
	// Get the main window from the glade file
	mainWindow, err := softInvoice.helper.GetApplicationWindow("main_window")
	errorCheck(err)

	// Set up main window
	mainWindow.SetApplication(app)
	mainWindow.SetTitle("Window")
	mainWindow.SetDefaultSize(800, 600)

	// Hook up the destroy event
	mainWindow.Connect("destroy", func() {
		closeApplication(softInvoice)
	})

	// Get the new invoice button
	button, err := softInvoice.helper.GetToolButton("newinvoice_button")
	errorCheck(err)

	// Hook up the clicked event for the new invoice button
	button.Connect("clicked", func() {
		openInvoiceDialog(softInvoice)
	})

	// Show the main window
	mainWindow.ShowAll()
}

func closeApplication(softInvoice *SoftInvoice) {
	// Destroy the invoice window if it has been created
	if softInvoice.invoiceWindow != nil {
		softInvoice.invoiceWindow.Destroy()
	}
	// Close the database
	softInvoice.database.CloseDatabase()
}

func openInvoiceDialog(softInvoice *SoftInvoice) {
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

// Create a new SoftInvoice object
func NewSoftInvoice(app *gtk.Application) *SoftInvoice {
	softInvoice := new(SoftInvoice)
	softInvoice.database = new(database)
	softInvoice.application = app
	return softInvoice
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
