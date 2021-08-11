package main

import (
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"os"
)

const (
	ApplicationId string = "se.softteam.invoice"
	ApplicationFlags = glib.APPLICATION_FLAGS_NONE
)

var softInvoice *SoftInvoice

func main() {
	// Parse command line arguments
	useTestDatabasePointer := flag.Bool("test",false, "Use testing database")
	flag.Parse()

	// Create a new application
	app, err := gtk.ApplicationNew(ApplicationId, ApplicationFlags)
	errorCheck(err)

	// Create the SoftInvoice application object
	softInvoice = NewSoftInvoice(app, *useTestDatabasePointer)
	softInvoice.mainForm = NewMainForm()

	// Hook up the activate event handler
	_ = app.Connect("activate", softInvoice.mainForm.OpenMainForm)

	// Start the application (and exit when it is done)
	os.Exit(app.Run(nil))
}

func errorCheck(err error) {
	if err != nil {
		panic(err)
	}
}
