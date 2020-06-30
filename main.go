package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"os"
)

const ApplicationId string = "se.softteam.invoice"
const ApplicationFlags glib.ApplicationFlags = glib.APPLICATION_FLAGS_NONE

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
