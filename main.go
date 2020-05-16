package main

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	gtkhelper "github.com/hultan/softteam/gtk"
	"os"
)

func main() {
	app, err:=gtk.ApplicationNew("se.softteam.invoice",glib.APPLICATION_FLAGS_NONE)
	errorCheck(err)
	app.Connect("activate", activate)
	status := app.Run(os.Args)

	os.Exit(status)
}

func activate(app *gtk.Application) {
	gtk.Init(&os.Args)
	builder, err := gtk.BuilderNewFromFile("resources/main.glade")
	errorCheck(err)

	var helper *gtkhelper.GtkHelper = gtkhelper.GtkHelperNew(builder)
	window, err:=helper.GetApplicationWindow("main_window")
	errorCheck(err)
	window.SetApplication(app)

	window.SetTitle("Window")
	window.SetDefaultSize(800,600)

	button, err := helper.GetToolButton("newinvoice_button")
	errorCheck(err)
	button.Connect("clicked", func() {
		openInvoiceDialog(app, helper)
	})

	window.ShowAll()
}

func openInvoiceDialog(app *gtk.Application, helper *gtkhelper.GtkHelper) {
	window, err := helper.GetApplicationWindow("invoice_window")
	errorCheck(err)
	window.SetApplication(app)
	window.SetModal(true)
	window.SetKeepAbove(true)
	button, err := helper.GetButton("cancel_button")
	errorCheck(err)
	button.Connect("clicked", func() {
		window.Close()
		window.Destroy()
	})

	window.Show()
}

func errorCheck(err error) {
	if err!=nil {
		panic(err)
	}
}