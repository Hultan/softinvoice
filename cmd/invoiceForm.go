package main

import (
	"database/sql"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/hultan/softteam-invoice/internal/database"
	"github.com/hultan/softteam/messagebox"
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
	invoiceRowTreeview    *gtk.TreeView
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

func (i *InvoiceForm) OpenInvoiceForm(softInvoice *SoftInvoice, reloadListCallback ReloadListCallback) {
	var err error
	i.reloadListCallback = reloadListCallback

	// Check if it is the first time we open the invoice window
	if softInvoice.invoiceForm.window == nil {
		// Get the invoice window from glade
		window, err := softInvoice.helper.GetWindow("invoice_window")
		errorCheck(err)

		// Save a pointer to the invoice window
		softInvoice.invoiceForm.window = window

		// Set up the invoice window
		window.SetApplication(softInvoice.application)
		window.HideOnDelete()
		window.SetModal(true)
		window.SetKeepAbove(true)
		window.SetPosition(gtk.WIN_POS_CENTER_ALWAYS)

		// Hook up the hide event
		_,err = window.Connect("hide", func() {
			i.CloseInvoiceWindow(softInvoice)
		})
		errorCheck(err)

		// Get the cancel button
		cancelButton, err := softInvoice.helper.GetButton("cancel_button")
		errorCheck(err)

		// Hook up the clicked event for the cancel button
		_,err = cancelButton.Connect("clicked", func() {
			window.Hide()
		})
		errorCheck(err)

		// Get the save button
		saveButton, err := softInvoice.helper.GetButton("save_button")
		errorCheck(err)

		// Hook up the clicked event for the save button
		_, err = saveButton.Connect("clicked", func() {
			isSaving = true
			window.Hide()
		})
		errorCheck(err)

		// Get the add button
		addButton, err := softInvoice.helper.GetButton("addrow_button")
		errorCheck(err)

		// Hook up the clicked event for the add row button
		_, err = addButton.Connect("clicked", func() {
			softInvoice.invoiceRowForm.OpenInvoiceRowForm(softInvoice, i.OnInvoiceRowAdded)
		})
		errorCheck(err)

		// Get row treeview
		treeview, err := softInvoice.helper.GetTreeView("invoicerow_treeview")
		errorCheck(err)
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
	nameEntry, err := softInvoice.helper.GetEntry("name_entry")
	if err != nil {
		fmt.Println("Failed to get name entry : ", err.Error())
	}
	i.nameEntry = nameEntry

	// Get address entry
	addressEntry, err := softInvoice.helper.GetEntry("address_entry")
	if err != nil {
		fmt.Println("Failed to get address entry : ", err.Error())
	}
	i.addressEntry = addressEntry

	// Get postal address entry
	postalAddressEntry, err := softInvoice.helper.GetEntry("postaladdress_entry")
	if err != nil {
		fmt.Println("Failed to get postal address entry : ", err.Error())
	}
	i.postalAddressEntry = postalAddressEntry

	// Get payday entry
	paydayEntry, err := softInvoice.helper.GetEntry("payday_entry")
	if err != nil {
		fmt.Println("Failed to get payday entry : ", err.Error())
	}
	i.paydayEntry = paydayEntry

	// Get your reference entry
	yourReferenceEntry, err := softInvoice.helper.GetEntry("yourreference_entry")
	if err != nil {
		fmt.Println("Failed to get your reference entry : ", err.Error())
	}
	i.yourReferenceEntry = yourReferenceEntry

	// Get invoice number entry
	invoiceNumberEntry, err := softInvoice.helper.GetEntry("invoicenumber_entry")
	if err != nil {
		fmt.Println("Failed to get invoice number entry : ", err.Error())
	}
	i.invoiceNumberEntry = invoiceNumberEntry

	// Get invoice date entry
	invoiceDateEntry, err := softInvoice.helper.GetEntry("invoicedate_entry")
	if err != nil {
		fmt.Println("Failed to get invoice date entry : ", err.Error())
	}
	i.invoiceDateEntry = invoiceDateEntry

	// Get calendar entry
	calendar, err := softInvoice.helper.GetCalendar("calendar")
	errorCheck(err)
	i.calendar = calendar
	_, err = calendar.Connect("day-selected", i.OnCalendarDateChanged)
	errorCheck(err)
}

func (i *InvoiceForm) SetupCustomerCombo(softInvoice *SoftInvoice) {
	// Get customer combo
	customerCombo, err := softInvoice.helper.GetComboBox("customer_combo")
	if err != nil {
		fmt.Println("Failed to get customer combobox : ", err.Error())
	}
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
		if iter!=nil {
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
	_, err = customerCombo.Connect("changed", i.OnCustomerChange)
	errorCheck(err)
}

//
// Signal handlers
//

func (i *InvoiceForm) OnCustomerChange(customerCombo *gtk.ComboBox) {
	// Get the id of the selected row
	iter, _ := customerCombo.GetActiveIter()
	model, _ := customerCombo.GetModel()
	idValue, _ := model.GetValue(iter, 0)
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
		messagebox.NewMessageBox("Missing customer...", "The invoice is missing a customer!", i.window)
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
		strconv.FormatFloat(float64(invoiceRow.Price),'f',2,32),
		strconv.FormatFloat(float64(invoiceRow.Amount),'f',2,32),
		strconv.FormatFloat(float64(invoiceRow.Total),'f',2,32),
	})

	i.invoiceRowTreeview.SetModel(i.rowListStore)
}
