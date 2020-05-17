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

type SoftInvoice struct {
	application *gtk.Application
	helper *gtkhelper.GtkHelper
	database *database
}

func main() {
	app, err:=gtk.ApplicationNew(ApplicationId,ApplicationFlags)
	errorCheck(err)

	softInvoice:=NewSoftInvoice(app)

	app.Connect("activate", activate, softInvoice)
	status := app.Run(os.Args)

	softInvoice.database.CloseDatabase()
	os.Exit(status)
}

func activate(app *gtk.Application, softInvoice *SoftInvoice) {
	gtk.Init(&os.Args)

	var helper *gtkhelper.GtkHelper = gtkhelper.GtkHelperNewFromFile("resources/main.glade")
	softInvoice.helper = helper

	window, err:=helper.GetApplicationWindow("main_window")
	errorCheck(err)
	window.SetApplication(app)

	window.SetTitle("Window")
	window.SetDefaultSize(800,600)

	button, err := helper.GetToolButton("newinvoice_button")
	errorCheck(err)
	button.Connect("clicked", func() {
		openInvoiceDialog(softInvoice)
	})

	window.ShowAll()
}

func openInvoiceDialog(softInvoice *SoftInvoice) {
	//window, err := softInvoice.helper.GetApplicationWindow("invoice_window")
	//errorCheck(err)
	//window.SetApplication(softInvoice.application)
	//window.SetModal(true)
	//window.SetKeepAbove(true)
	//softInvoice.application.AddWindow(window)

	invoices, err := softInvoice.database.GetAllInvoices()
	softErrorCheck(err)
	fmt.Println(invoices[10].CustomerName)
	fmt.Println(invoices[10].rows[0].Price)

	invoice, err := softInvoice.database.GetInvoiceByNumber("1245")
	softErrorCheck(err)
	fmt.Println(invoice.CustomerName)
	fmt.Println(invoice.rows[0].Price)

	customers, err := softInvoice.database.GetAllCustomers()
	softErrorCheck(err)
	fmt.Println(customers[1].FancyName)

	customer, err := softInvoice.database.GetCustomerByNumber("1021")
	softErrorCheck(err)
	fmt.Println(customer.FancyName)

	// Should fail
	_, err = softInvoice.database.GetCustomerByNumber("1234")
	softErrorCheck(err)


	products, err := softInvoice.database.GetAllProducts()
	softErrorCheck(err)
	fmt.Println(products[1].Name)

	product, err := softInvoice.database.GetProductByNumber("NOVA")
	softErrorCheck(err)
	fmt.Println(product.Name)

	// Should fail
	_, err = softInvoice.database.GetProductByNumber("FAIL")
	softErrorCheck(err)

	//button, err := softInvoice.helper.GetButton("cancel_button")
	//errorCheck(err)
	//button.Connect("clicked", func() {
	//	window.Hide()
	//})

}

func errorCheck(err error) {
	if err!=nil {
		panic(err)
	}
}

func softErrorCheck(err error) {
	if err!=nil {
		fmt.Println(err.Error())
	}
}

func NewSoftInvoice(app *gtk.Application) *SoftInvoice {
	softInvoice:=new(SoftInvoice)
	softInvoice.database = new(database)
	softInvoice.application = app
	return softInvoice
}