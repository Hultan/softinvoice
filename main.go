package main

import (
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"os"
)

const (
	ApplicationId string = "se.softteam.invoice"
	ApplicationFlags glib.ApplicationFlags = glib.APPLICATION_FLAGS_NONE
)

var (
	useTestDatabase = false
)

func main() {
	useTestDatabase = *flag.Bool("t",false, "Use testing database")
	flag.Parse()

	// Check command line arguments
	//if len(os.Args) > 1 {
	//	if strings.HasPrefix(os.Args[1],"-t") {
	//		useTestDatabase = true
	//	}
	//}

	// Create a new application
	app, err := gtk.ApplicationNew(ApplicationId, ApplicationFlags)
	errorCheck(err)

	// Create the SoftInvoice application object
	softInvoice := NewSoftInvoice(app)
	softInvoice.mainForm = NewMainForm()

	// Hook up the activate event handler
	app.Connect("activate", softInvoice.mainForm.OpenMainForm, softInvoice)

	// Start the application (and exit when it is done)
	os.Exit(app.Run(nil))
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
