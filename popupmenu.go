package main

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/hultan/softteam-invoice/database"
	"strconv"
)

type PopupMenu struct {
	parent *MainForm
}

func NewPopupMenu(softInvoice *SoftInvoice, mainWindow *MainForm) *PopupMenu {
	menu := new(PopupMenu)

	menu.parent = mainWindow

	popup, err := softInvoice.helper.GetMenu("popup_menu")
	errorCheck(err)

	preview, err := softInvoice.helper.GetMenuItem("popupMenuItemPreview")
	errorCheck(err)

	pdf, err := softInvoice.helper.GetMenuItem("popupMenuItemSaveAsPDF")
	errorCheck(err)

	mainWindow.treeView.Connect("button-release-event", func(treeview *gtk.TreeView, event *gdk.Event) {
		buttonEvent := gdk.EventButtonNewFromEvent(event)
		if buttonEvent.Button() == 3 { // 3 == Mouse right button!?
			popup.PopupAtPointer(event)
		}
	})

	preview.Connect("activate", func() {
		invoice := menu.getSelectedInvoice(mainWindow.treeView)
		if invoice == nil {
			return
		}

		softInvoice.previewWForm.OpenPreviewForm(softInvoice, invoice)
	})

	pdf.Connect("activate", func() {
		invoice := menu.getSelectedInvoice(mainWindow.treeView)
		if invoice == nil {
			return
		}

		dialog, err := gtk.FileChooserDialogNewWith2Buttons("Save PDF as...", mainWindow.window,gtk.FILE_CHOOSER_ACTION_SAVE,
			"Cancel", gtk.RESPONSE_CANCEL, "Save", gtk.RESPONSE_ACCEPT)
		errorCheck(err)

		response := dialog.Run()
		if response==gtk.RESPONSE_ACCEPT {
			creator := NewInvoiceCreator(invoice)
			creator.CreatePDF(dialog.GetFilename())
		}

		dialog.Destroy()
	})
	return menu
}

func (p *PopupMenu) getSelectedInvoice(treeView *gtk.TreeView) *database.Invoice {
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
		for _, invoice := range p.parent.invoices {
			if invoice.Number == invoiceNumber {
				return &invoice
			}
		}
		return nil
	}

	return nil
}
