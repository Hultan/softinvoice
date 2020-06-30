package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	database "github.com/hultan/softteam-invoice/database"
	gtkhelper "github.com/hultan/softteam/gtk"
	"os"
)

const ApplicationId string = "se.softteam.invoice"
const ApplicationFlags glib.ApplicationFlags = glib.APPLICATION_FLAGS_NONE

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

func main() {
	// Create a new application
	app, err := gtk.ApplicationNew(ApplicationId, ApplicationFlags)
	errorCheck(err)

	// Create the SoftInvoice application object
	softInvoice := NewSoftInvoice(app)
	softInvoice.mainForm = NewMainForm()

	// Hook up the activate event handler
	app.Connect("activate", softInvoice.mainForm.OpenMainForm, softInvoice)

	// Start the application
	status := app.Run(os.Args)

	// Exit
	os.Exit(status)
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
