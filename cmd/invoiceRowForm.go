package main

import (
	"fmt"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/hultan/softteam-invoice/internal/database"
	"strconv"
)

type SaveCallback func(*database.InvoiceRow)

type InvoiceRowForm struct {
	window *gtk.Window

	productCombo *gtk.ComboBox
	textEntry    *gtk.Entry
	nameEntry    *gtk.Entry
	priceEntry   *gtk.Entry
	amountEntry  *gtk.Entry

	products []database.Product

	SaveCallback SaveCallback
}

var isSavingRow = false

func NewInvoiceRowForm() *InvoiceRowForm {
	invoiceRowForm := new(InvoiceRowForm)
	return invoiceRowForm
}

func (i *InvoiceRowForm) OpenInvoiceRowForm(softInvoice *SoftInvoice, saveCallback SaveCallback) {
	i.SaveCallback = saveCallback

	// Check if it is the first time we open the invoice row window
	if softInvoice.invoiceRowForm.window == nil {
		// Get the invoice window from glade
		window, err := softInvoice.helper.GetWindow("invoicerow_window")
		errorCheck(err)

		// Save a pointer to the invoice window
		softInvoice.invoiceRowForm.window = window

		// Set up the invoice window
		window.SetApplication(softInvoice.application)
		window.HideOnDelete()
		window.SetModal(true)
		window.SetKeepAbove(true)
		window.SetPosition(gtk.WIN_POS_CENTER_ALWAYS)

		// Hook up the hide event
		_, err = window.Connect("hide", func() {
			i.CloseInvoiceRowWindow()
		})
		errorCheck(err)

		// Get the cancel button
		cancelButton, err := softInvoice.helper.GetButton("productcancel_button")
		errorCheck(err)

		// Hook up the clicked event for the cancel button
		_, err = cancelButton.Connect("clicked", func() {
			window.Hide()
		})
		errorCheck(err)

		// Get the save button
		saveButton, err := softInvoice.helper.GetButton("productsave_button")
		errorCheck(err)

		// Hook up the clicked event for the save button
		_, err = saveButton.Connect("clicked", func() {
			isSavingRow = true
			window.Hide()
		})
		errorCheck(err)

		// Setup window
		i.SetupWindow(softInvoice)
		i.SetupProductCombo(softInvoice)
	}

	// Set default values
	i.productCombo.SetActive(0)
	i.amountEntry.SetText("1")

	// Show the window
	softInvoice.invoiceRowForm.window.ShowAll()
}

func (i *InvoiceRowForm) CloseInvoiceRowWindow() {
	if isSavingRow {
		isSavingRow = false
		row := i.SaveInvoiceRow()
		i.SaveCallback(row)
	}
}

//
//Setup functions
//

func (i *InvoiceRowForm) SetupWindow(softInvoice *SoftInvoice) {
	// Get name entry
	nameEntry, err := softInvoice.helper.GetEntry("productname_entry")
	if err != nil {
		fmt.Println("Failed to get name entry : ", err.Error())
	}
	i.nameEntry = nameEntry

	// Get text entry
	textEntry, err := softInvoice.helper.GetEntry("producttext_entry")
	if err != nil {
		fmt.Println("Failed to get text entry : ", err.Error())
	}
	i.textEntry = textEntry

	// Get price entry
	priceEntry, err := softInvoice.helper.GetEntry("productprice_entry")
	if err != nil {
		fmt.Println("Failed to get price entry : ", err.Error())
	}
	i.priceEntry = priceEntry

	// Get amount entry
	amountEntry, err := softInvoice.helper.GetEntry("productamount_entry")
	if err != nil {
		fmt.Println("Failed to get amount entry : ", err.Error())
	}
	i.amountEntry = amountEntry
}

func (i *InvoiceRowForm) SetupProductCombo(softInvoice *SoftInvoice) {
	// Get product combo
	productCombo, err := softInvoice.helper.GetComboBox("product_combo")
	if err != nil {
		fmt.Println("Failed to get product combobox : ", err.Error())
	}
	i.productCombo = productCombo

	// Get all products from the database
	products, err := softInvoice.database.GetAllProducts()
	if err != nil {
		fmt.Println("Failed to load products : ", err.Error())
	}
	i.products = products

	// Add product to a list store
	productStore, err := gtk.ListStoreNew(glib.TYPE_INT, glib.TYPE_STRING, glib.TYPE_STRING)
	errorCheck(err)
	for _, value := range products {
		iter := productStore.Append()
		// Ignore errors here
		_ = productStore.Set(iter, []int{0, 1, 2}, []interface{}{value.Id, value.Number, value.Name})
	}

	// Setup combo and renderer
	productCombo.SetModel(productStore)
	numberRenderer, _ := gtk.CellRendererTextNew()
	productCombo.PackStart(numberRenderer, true)
	productCombo.AddAttribute(numberRenderer, "text", 1)

	nameRenderer, _ := gtk.CellRendererTextNew()
	productCombo.PackStart(nameRenderer, true)
	productCombo.AddAttribute(nameRenderer, "text", 2)

	_, err = productCombo.Connect("changed", i.OnProductChange)
	errorCheck(err)
}

//
// Signal handlers
//

func (i *InvoiceRowForm) OnProductChange(customerCombo *gtk.ComboBox) {
	// Get the id of the selected product
	iter, _ := customerCombo.GetActiveIter()
	model, _ := customerCombo.GetModel()
	idValue, _ := model.GetValue(iter, 0)
	id, _ := idValue.GoValue()

	// Find the selected product
	var foundProduct database.Product
	var found = false

	for _, value := range i.products {
		if value.Id == id.(int) {
			foundProduct = value
			found = true
			break
		}
	}

	if !found {
		panic("Customer not found!")
	}

	// Set some product related fields
	i.nameEntry.SetText(foundProduct.Name)
	i.textEntry.SetText(foundProduct.Text)
	i.priceEntry.SetText(fmt.Sprintf("%.0f", foundProduct.Price))
}

//
// Save function
//

func (i *InvoiceRowForm) SaveInvoiceRow() *database.InvoiceRow {
	var row database.InvoiceRow

	// Get text and name
	row.Text, _ = i.textEntry.GetText()
	row.Name, _ = i.nameEntry.GetText()

	// Get and parse the price field
	priceString, _ := i.priceEntry.GetText()
	price, _ := strconv.ParseFloat(priceString, 32)
	row.Price = float32(price)

	// Get and parse the amount field
	amountString, _ := i.amountEntry.GetText()
	amount, _ := strconv.ParseFloat(amountString, 32)
	row.Amount = float32(amount)

	// Calculate the row total (excl VAT)
	row.Total = float32(amount * price)

	return &row
}
