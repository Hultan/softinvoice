package main

import (
	"fmt"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/hultan/softinvoice/internal/database"
	"strconv"
	"strings"
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
		window := softInvoice.builder.GetObject("invoicerow_window").(*gtk.Window)

		// Save a pointer to the invoice window
		softInvoice.invoiceRowForm.window = window

		// Set up the invoice window
		window.SetApplication(softInvoice.application)
		window.HideOnDelete()
		window.SetModal(true)
		window.SetKeepAbove(true)
		window.SetPosition(gtk.WIN_POS_CENTER_ALWAYS)
		window.SetTransientFor(softInvoice.invoiceForm.window)

		// Hook up the hide event
		_ = window.Connect("hide", func() {
			i.CloseInvoiceRowWindow()
		})

		// Get the cancel button
		cancelButton := softInvoice.builder.GetObject("productcancel_button").(*gtk.Button)

		// Hook up the clicked event for the cancel button
		_ = cancelButton.Connect("clicked", func() {
			window.Hide()
		})

		// Get the save button
		saveButton := softInvoice.builder.GetObject("productsave_button").(*gtk.Button)

		// Hook up the clicked event for the save button
		_ = saveButton.Connect("clicked", func() {
			isSavingRow = true
			window.Hide()
		})

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
	nameEntry := softInvoice.builder.GetObject("productname_entry").(*gtk.Entry)
	i.nameEntry = nameEntry

	// Get text entry
	textEntry := softInvoice.builder.GetObject("producttext_entry").(*gtk.Entry)
	i.textEntry = textEntry

	// Get price entry
	priceEntry := softInvoice.builder.GetObject("productprice_entry").(*gtk.Entry)
	i.priceEntry = priceEntry

	// Get amount entry
	amountEntry := softInvoice.builder.GetObject("productamount_entry").(*gtk.Entry)
	i.amountEntry = amountEntry
}

func (i *InvoiceRowForm) SetupProductCombo(softInvoice *SoftInvoice) {
	// Get product combo
	productCombo := softInvoice.builder.GetObject("product_combo").(*gtk.ComboBox)
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

	_ = productCombo.Connect("changed", i.OnProductChange)
}

//
// Signal handlers
//

func (i *InvoiceRowForm) OnProductChange(customerCombo *gtk.ComboBox) {
	// Get the id of the selected product
	iter, _ := customerCombo.GetActiveIter()
	model, _ := customerCombo.GetModel()
	idValue, _ := model.(*gtk.TreeModel).GetValue(iter, 0)
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
	amountString = strings.Replace(amountString,",",".", 1)
	amount, _ := strconv.ParseFloat(amountString, 32)
	row.Amount = float32(amount)

	// Calculate the row total (excl VAT)
	row.Total = float32(amount * price)

	return &row
}
