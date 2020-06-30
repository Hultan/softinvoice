package main

import (
	"database/sql"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/hultan/softteam-invoice/database"
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
		window.Connect("hide", func() {
			i.CloseInvoiceWindow(softInvoice)
		})

		// Get the cancel button
		cancelButton, err := softInvoice.helper.GetButton("cancel_button")
		errorCheck(err)

		// Hook up the clicked event for the cancel button
		cancelButton.Connect("clicked", func() {
			window.Hide()
		})

		// Get the save button
		saveButton, err := softInvoice.helper.GetButton("save_button")
		errorCheck(err)

		// Hook up the clicked event for the save button
		saveButton.Connect("clicked", func() {
			isSaving = true
			window.Hide()
		})

		// Get the add button
		addButton, err := softInvoice.helper.GetButton("addrow_button")
		errorCheck(err)

		// Hook up the clicked event for the add row button
		addButton.Connect("clicked", func() {
			softInvoice.invoiceRowForm.OpenInvoiceRowForm(softInvoice, i.invoiceRowAdded)
		})
	}

	// Setup window
	if i.customerCombo == nil {
		i.setupWindow(softInvoice)
		i.setupCustomerCombo(softInvoice)
	}

	// Set default values
	i.invoice = database.Invoice{}
	i.customer = database.Customer{}
	i.customerCombo.SetActive(0)
	currentTime := time.Now()
	i.calendar.SelectMonth(uint(currentTime.Month())-1, uint(currentTime.Year()))
	i.calendar.SelectDay(uint(currentTime.Day()))

	// Show the window
	softInvoice.invoiceForm.window.ShowAll()
}

func (i *InvoiceForm) CloseInvoiceWindow(softInvoice *SoftInvoice) {
	if isSaving {
		isSaving = false
		i.saveInvoice(softInvoice)
		i.reloadListCallback(softInvoice)
	}
}

func (i *InvoiceForm) setupWindow(softInvoice *SoftInvoice) {
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
	if err != nil {
		fmt.Println("Failed to get calendar : ", err.Error())
	}
	i.calendar = calendar

	calendar.Connect("day-selected", func() {
		year, month, day := i.calendar.GetDate()
		date := fmt.Sprintf("%d-%.2d-%.2d", year, month+1, day)
		invoiceDateEntry.SetText(date)
	})
}

func (i *InvoiceForm) setupCustomerCombo(softInvoice *SoftInvoice) {
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

	// Add customers to a list store
	customerStore, err := gtk.ListStoreNew(glib.TYPE_INT, glib.TYPE_STRING, glib.TYPE_STRING)
	for _, value := range customers {
		iter := customerStore.Append()
		customerStore.Set(iter, []int{0, 1, 2}, []interface{}{value.Id, value.Number, value.Name})
	}

	// Setup combo and renderer
	customerCombo.SetModel(customerStore)
	numberRenderer, _ := gtk.CellRendererTextNew()
	customerCombo.PackStart(numberRenderer, true)
	customerCombo.AddAttribute(numberRenderer, "text", 1)

	nameRenderer, _ := gtk.CellRendererTextNew()
	customerCombo.PackStart(nameRenderer, true)
	customerCombo.AddAttribute(nameRenderer, "text", 2)

	customerCombo.Connect("changed", i.onCustomerChange, softInvoice)
}

func (i *InvoiceForm) onCustomerChange(customerCombo *gtk.ComboBox, softInvoice *SoftInvoice) {
	iter, _ := customerCombo.GetActiveIter()
	model, _ := customerCombo.GetModel()
	idValue, _ := model.GetValue(iter, 0)
	id, _ := idValue.GoValue()

	var foundCustomer database.Customer
	var found bool = false

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

	i.nameEntry.SetText(foundCustomer.Name)
	i.addressEntry.SetText(foundCustomer.Address)
	i.postalAddressEntry.SetText(foundCustomer.PostalAddress)
	i.paydayEntry.SetText(strconv.Itoa(foundCustomer.PayDay))
	i.yourReferenceEntry.SetText(foundCustomer.Reference)
	next, err := softInvoice.database.GetNextInvoiceNumber()
	if err != nil {
		panic(fmt.Sprintf("Failed to get next invoice number : %s", err.Error()))
	}
	i.invoiceNumberEntry.SetText(strconv.Itoa(next))
}

func (i *InvoiceForm) saveInvoice(softInvoice *SoftInvoice) bool {
	if i.customer.Number=="" {
		messagebox.NewMessageBox("Missing customer...","The invoice is missing a customer!", i.window)
		panic("missing customer")
	}

	number, _ := i.invoiceNumberEntry.GetText()
	i.invoice.Number, _ = strconv.Atoi(number)

	dateString, _ := i.invoiceDateEntry.GetText()
	date, _ := time.Parse(constDateLayout, dateString)
	i.invoice.Date = date

	payDayString, _ := i.paydayEntry.GetText()
	payDay, _ := strconv.Atoi(payDayString)
	i.invoice.PayDay = payDay

	dueDate := date.AddDate(0, 0, payDay)
	i.invoice.DueDate = dueDate

	i.invoice.CustomerNumber = i.customer.Number
	i.invoice.CustomerName = i.customer.Name
	i.invoice.CustomerAddress = i.customer.Address
	i.invoice.CustomerPostalAddress = i.customer.PostalAddress
	i.invoice.CustomerReference = i.customer.Reference

	i.invoice.Credit = false
	i.invoice.CreditInvoiceNumber = sql.NullInt32 {
		Int32: 0,
		Valid: false,
	}

	i.invoice.ReadOnly = false

	var amountWithoutVAT float32
	for _, value := range i.invoice.Rows {
		amountWithoutVAT += value.Total
	}
	i.invoice.Amount = amountWithoutVAT * 1.25

	// Pretty print (spew) the invoice
	spew.Dump(i.invoice)

	// Save invoice
	err := softInvoice.database.InsertInvoice(&i.invoice)
	if err!=nil {
		fmt.Println(err.Error())
	}

	return true
}

func (i *InvoiceForm) invoiceRowAdded(row *database.InvoiceRow) {
	i.invoice.Rows = append(i.invoice.Rows, *row)
}
