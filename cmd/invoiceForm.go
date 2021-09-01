package main

import (
	"database/sql"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/hultan/softinvoice/internal/database"
	"strconv"
	"time"
)

type ReloadListCallback func(*SoftInvoice) error

type InvoiceForm struct {
	window *gtk.Window

	customerCombo      *gtk.ComboBox
	nameEntry          *gtk.Entry
	addressEntry       *gtk.Entry
	postalAddressEntry *gtk.Entry
	paydayEntry        *gtk.Entry
	yourReferenceEntry *gtk.Entry
	invoiceNumberEntry *gtk.Entry
	invoiceDateEntry   *gtk.Entry
	calendar           *gtk.Calendar
	invoiceRowTreeview *gtk.TreeView
	rowListStore       *gtk.ListStore

	customer  database.Customer
	customers []database.Customer
	invoice   database.Invoice

	reloadListCallback ReloadListCallback
}

var isSaving = false

func NewInvoiceForm() *InvoiceForm {
	invoiceForm := new(InvoiceForm)
	return invoiceForm
}

func (i *InvoiceForm) OpenInvoiceForm(reloadListCallback ReloadListCallback) {
	var err error
	i.reloadListCallback = reloadListCallback

	// Check if it is the first time we open the invoice window
	if softInvoice.invoiceForm.window == nil {
		// Get the invoice window from glade
		window := softInvoice.builder.GetObject("invoice_window").(*gtk.Window)

		// Save a pointer to the invoice window
		softInvoice.invoiceForm.window = window

		// Set up the invoice window
		window.SetApplication(softInvoice.application)
		window.HideOnDelete()
		window.SetModal(true)
		window.SetKeepAbove(true)
		window.SetPosition(gtk.WIN_POS_CENTER_ALWAYS)

		// Hook up the hide event
		_ = window.Connect("hide", func() {
			i.CloseInvoiceWindow(softInvoice)
		})

		// Get the cancel button
		cancelButton := softInvoice.builder.GetObject("cancel_button").(*gtk.Button)

		// Hook up the clicked event for the cancel button
		_ = cancelButton.Connect("clicked", func() {
			window.Hide()
		})

		// Get the save button
		saveButton := softInvoice.builder.GetObject("save_button").(*gtk.Button)

		// Hook up the clicked event for the save button
		_ = saveButton.Connect("clicked", func() {
			isSaving = true
			window.Hide()
		})

		// Get the add button
		addButton := softInvoice.builder.GetObject("addrow_button").(*gtk.Button)

		// Hook up the clicked event for the add row button
		_ = addButton.Connect("clicked", func() {
			softInvoice.invoiceRowForm.OpenInvoiceRowForm(softInvoice, i.OnInvoiceRowAdded)
		})

		// Get row treeview
		treeview := softInvoice.builder.GetObject("invoicerow_treeview").(*gtk.TreeView)
		i.invoiceRowTreeview = treeview
	}

	// Setup window
	if i.customerCombo == nil {
		i.SetupWindow(softInvoice)
		i.SetupCustomerCombo(softInvoice)
	}

	// Create row list store
	i.rowListStore, err = gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING)
	errorCheck(err)
	i.invoiceRowTreeview.SetModel(i.rowListStore)

	// Set default values
	i.invoice = database.Invoice{}
	i.customer = database.Customer{}
	i.customerCombo.SetActive(0)
	currentTime := time.Now()
	i.calendar.SelectMonth(uint(currentTime.Month())-1, uint(currentTime.Year()))
	i.calendar.SelectDay(uint(currentTime.Day()))
	next, err := softInvoice.database.GetNextInvoiceNumber()
	if err != nil {
		panic(fmt.Sprintf("Failed to get next invoice number : %s", err.Error()))
	}
	i.invoiceNumberEntry.SetText(strconv.Itoa(next))

	// Show the window
	softInvoice.invoiceForm.window.ShowAll()
}

func (i *InvoiceForm) CloseInvoiceWindow(softInvoice *SoftInvoice) {
	if isSaving {
		isSaving = false
		// Save the new invoice
		i.SaveInvoice(softInvoice)
		// Make sure the main form reloads the list of invoices
		err := i.reloadListCallback(softInvoice)
		errorCheck(err)
	}
}

//
// Setup functions
//

func (i *InvoiceForm) SetupWindow(softInvoice *SoftInvoice) {
	// Get name entry
	nameEntry := softInvoice.builder.GetObject("name_entry").(*gtk.Entry)
	i.nameEntry = nameEntry

	// Get address entry
	addressEntry := softInvoice.builder.GetObject("address_entry").(*gtk.Entry)
	i.addressEntry = addressEntry

	// Get postal address entry
	postalAddressEntry := softInvoice.builder.GetObject("postaladdress_entry").(*gtk.Entry)
	i.postalAddressEntry = postalAddressEntry

	// Get payday entry
	paydayEntry := softInvoice.builder.GetObject("payday_entry").(*gtk.Entry)
	i.paydayEntry = paydayEntry

	// Get your reference entry
	yourReferenceEntry := softInvoice.builder.GetObject("yourreference_entry").(*gtk.Entry)
	i.yourReferenceEntry = yourReferenceEntry

	// Get invoice number entry
	invoiceNumberEntry := softInvoice.builder.GetObject("invoicenumber_entry").(*gtk.Entry)
	i.invoiceNumberEntry = invoiceNumberEntry

	// Get invoice date entry
	invoiceDateEntry := softInvoice.builder.GetObject("invoicedate_entry").(*gtk.Entry)
	i.invoiceDateEntry = invoiceDateEntry

	// Get calendar entry
	calendar := softInvoice.builder.GetObject("calendar").(*gtk.Calendar)
	i.calendar = calendar
	_ = calendar.Connect("day-selected", i.OnCalendarDateChanged)
}

func (i *InvoiceForm) SetupCustomerCombo(softInvoice *SoftInvoice) {
	// Get customer combo
	customerCombo := softInvoice.builder.GetObject("customer_combo").(*gtk.ComboBox)
	i.customerCombo = customerCombo

	// Get all customers from the database
	customers, err := softInvoice.database.GetAllCustomers()
	if err != nil {
		fmt.Println("Failed to load customers : ", err.Error())
	}
	i.customers = customers

	var iter *gtk.TreeIter
	// Add customers to a list store
	customerStore, err := gtk.ListStoreNew(glib.TYPE_INT, glib.TYPE_STRING, glib.TYPE_STRING)
	errorCheck(err)
	for _, value := range customers {
		iter = customerStore.Append()
		if iter != nil {
			// Don't add customer on error
			_ = customerStore.Set(iter, []int{0, 1, 2}, []interface{}{value.Id, value.Number, value.Name})
		}
	}

	// Setup combo and renderer
	customerCombo.SetModel(customerStore)
	numberRenderer, _ := gtk.CellRendererTextNew()
	customerCombo.PackStart(numberRenderer, true)
	customerCombo.AddAttribute(numberRenderer, "text", 1)

	nameRenderer, _ := gtk.CellRendererTextNew()
	customerCombo.PackStart(nameRenderer, true)
	customerCombo.AddAttribute(nameRenderer, "text", 2)

	// Hook up customer change signal
	_ = customerCombo.Connect("changed", i.OnCustomerChange)
}

//
// Signal handlers
//

func (i *InvoiceForm) OnCustomerChange(customerCombo *gtk.ComboBox) {
	// Get the id of the selected row
	iter, _ := customerCombo.GetActiveIter()
	model, _ := customerCombo.GetModel()
	idValue, _ := model.(*gtk.TreeModel).GetValue(iter, 0)
	id, _ := idValue.GoValue()

	// Loop through customers and find the selected one
	var foundCustomer database.Customer
	var found = false

	for _, value := range i.customers {
		if value.Id == id.(int) {
			foundCustomer = value
			found = true
			break
		}
	}

	if !found {
		panic("Customer not found!")
	}
	i.customer = foundCustomer

	// Set some customer related fields
	i.nameEntry.SetText(foundCustomer.Name)
	i.addressEntry.SetText(foundCustomer.Address)
	i.postalAddressEntry.SetText(foundCustomer.PostalAddress)
	i.paydayEntry.SetText(strconv.Itoa(foundCustomer.PayDay))
	i.yourReferenceEntry.SetText(foundCustomer.Reference)
}

func (i *InvoiceForm) OnInvoiceRowAdded(row *database.InvoiceRow) {
	i.invoice.Rows = append(i.invoice.Rows, *row)
	i.AddInvoiceRow(row)
}

func (i *InvoiceForm) OnCalendarDateChanged() {
	year, month, day := i.calendar.GetDate()
	date := fmt.Sprintf("%d-%.2d-%.2d", year, month+1, day)
	i.invoiceDateEntry.SetText(date)
}

//
// Save invoice function
//

func (i *InvoiceForm) SaveInvoice(softInvoice *SoftInvoice) bool {
	// Check that a customer has been selected
	if i.customer.Number == "" {
		messagebox := gtk.MessageDialogNew(i.window, gtk.DIALOG_MODAL, gtk.MESSAGE_ERROR, gtk.BUTTONS_OK, "The invoice is missing a customer!")
		messagebox.Show()
		panic("missing customer")
	}

	// Get invoice number
	number, _ := i.invoiceNumberEntry.GetText()
	i.invoice.Number, _ = strconv.Atoi(number)

	// Get and parse invoice date
	dateString, _ := i.invoiceDateEntry.GetText()
	date, _ := time.Parse(constDateLayout, dateString)
	i.invoice.Date = date

	// Get and parse paydays
	payDayString, _ := i.paydayEntry.GetText()
	payDay, _ := strconv.Atoi(payDayString)
	i.invoice.PayDay = payDay

	// Calculate due date
	dueDate := date.AddDate(0, 0, payDay)
	i.invoice.DueDate = dueDate

	// Set some customer related fields
	i.invoice.CustomerNumber = i.customer.Number
	i.invoice.CustomerName = i.customer.Name
	i.invoice.CustomerAddress = i.customer.Address
	i.invoice.CustomerPostalAddress = i.customer.PostalAddress
	i.invoice.CustomerReference = i.customer.Reference

	// Credit invoices (not done yet)
	i.invoice.Credit = false
	i.invoice.CreditInvoiceNumber = sql.NullInt32{
		Int32: 0,
		Valid: false,
	}

	// Handle read only flag (not done yet)
	i.invoice.ReadOnly = false

	// Calculate amounts
	var amountWithoutVAT float32
	for _, value := range i.invoice.Rows {
		amountWithoutVAT += value.Total
	}
	i.invoice.Amount = amountWithoutVAT * 1.25

	// Pretty print (spew) the invoice
	spew.Dump(i.invoice)

	// Save invoice
	err := softInvoice.database.InsertInvoice(&i.invoice)
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}

	return true
}

//
// Invoice rows
//

func (i *InvoiceForm) AddInvoiceRow(invoiceRow *database.InvoiceRow) {
	i.invoiceRowTreeview.SetModel(nil)

	iter := i.rowListStore.Append()
	// Ignore errors here
	_ = i.rowListStore.Set(iter, []int{0, 1, 2, 3, 4}, []interface{}{
		invoiceRow.Text,
		invoiceRow.Name,
		strconv.FormatFloat(float64(invoiceRow.Price), 'f', 2, 32),
		strconv.FormatFloat(float64(invoiceRow.Amount), 'f', 2, 32),
		strconv.FormatFloat(float64(invoiceRow.Total), 'f', 2, 32),
	})

	i.invoiceRowTreeview.SetModel(i.rowListStore)
}
