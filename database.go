package main

// https://github.com/jmoiron/sqlx	: Querying
// https://github.com/jmoiron/modl	: Inserts and updates
// Crap! Use gorm instead!

import (
	//	"github.com/jmoiron/sqlx"
	"github.com/jinzhu/gorm"
)

type database struct {
	db *gorm.DB
}

func (d *database) UpdateProduct(product *Product) error {
	db, err := d.getDatabase()
	if err != nil {
		return err
	}
	if result := db.Save(product); result.Error != nil {
		return result.Error
	}

	return nil
}

func (d *database) GetAllInvoices() ([]Invoice, error) {
	db, err := d.getDatabase()
	if err != nil {
		return nil, err
	}
	var invoices []Invoice
	if result := db.Find(&invoices); result.Error != nil {
		return nil, result.Error
	}

	for index, invoice := range invoices {
		rows, err := d.GetInvoiceRowsByInvoiceId(invoice.Id)
		if err != nil {
			return nil, err
		}
		invoices[index].rows = rows
	}

	return invoices, nil
}

func (d *database) GetInvoiceRowsByInvoiceId(id int) ([]InvoiceRow, error) {
	db, err := d.getDatabase()
	if err != nil {
		return nil, err
	}
	var invoiceRows []InvoiceRow
	if result := db.Where("invoiceid=?", id).Find(&invoiceRows); result.Error != nil {
		return nil, result.Error
	}

	return invoiceRows, nil
}

func (d *database) GetInvoiceByNumber(number string) (*Invoice, error) {
	db, err := d.getDatabase()
	if err != nil {
		return nil, err
	}
	invoice := Invoice{}
	if result := db.Where("number=?", number).First(&invoice); result.Error != nil {
		return nil, err
	}

	rows, err := d.GetInvoiceRowsByInvoiceId(invoice.Id)
	if err != nil {
		return nil, err
	}
	invoice.rows = rows

	return &invoice, nil
}

func (d *database) GetAllCustomers() ([]Customer, error) {
	db, err := d.getDatabase()
	if err != nil {
		return nil, err
	}
	var customers []Customer
	if result := db.Find(&customers); result.Error != nil {
		return nil, result.Error
	}

	return customers, nil
}

func (d *database) GetCustomerByNumber(number string) (*Customer, error) {
	db, err := d.getDatabase()
	if err != nil {
		return nil, err
	}
	customer := Customer{}
	if result := db.Where("number=?", number).First(&customer); result.Error != nil {
		return nil, result.Error
	}

	return &customer, nil
}

func (d *database) GetAllProducts() ([]Product, error) {
	db, err := d.getDatabase()
	if err != nil {
		return nil, err
	}
	var products []Product
	if result := db.Find(&products); result.Error != nil {
		return nil, result.Error
	}

	return products, nil
}

func (d *database) GetProductByNumber(number string) (*Product, error) {
	db, err := d.getDatabase()
	if err != nil {
		return nil, err
	}
	product := Product{}
	if result := db.Where("number=?", number).First(&product); result.Error != nil {
		return nil, result.Error
	}

	return &product, nil
}

func (d *database) getDatabase() (*gorm.DB, error) {
	if d.db == nil {
		db, err := gorm.Open("mysql", "per:KnaskimGjwQ6M!@tcp(192.168.1.3:3306)/softinvoice?parseTime=True")
		if err != nil {
			return nil, err
		}
		d.db = db
	}
	return d.db, nil
}

func (d *database) CloseDatabase() {
	if d.db == nil {
		return
	}
	d.db.Close()
	d.db = nil
	return
}
