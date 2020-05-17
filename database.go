package main

import (
	"github.com/jmoiron/sqlx"
)

type database struct {
	db *sqlx.DB
}

func (d *database) GetAllInvoices() ([]Invoice,error) {
	db, err := d.getDatabase()
	if err!= nil {
		return nil ,err
	}
	invoices := []Invoice{}
	err = db.Select(&invoices, "SELECT * FROM invoice")
	if err!=nil {
		return nil, err
	}

	for index, invoice := range invoices {
		rows, err :=d.GetInvoiceRowsByInvoiceId(invoice.Id)
		if err!=nil {
			return nil, err
		}
		invoices[index].rows = rows
	}

	return invoices, nil
}

func (d *database) GetInvoiceRowsByInvoiceId(id int) ([]InvoiceRow,error) {
	db, err := d.getDatabase()
	if err!= nil {
		return nil ,err
	}
	invoiceRows := []InvoiceRow{}
	err = db.Select(&invoiceRows, "SELECT * FROM invoicerow WHERE invoiceid=?", id)
	if err!= nil {
		return nil ,err
	}

	return invoiceRows, nil
}

func (d *database) GetInvoiceByNumber(number string) (*Invoice,error) {
	db, err := d.getDatabase()
	if err!= nil {
		return nil ,err
	}
	invoice := Invoice{}
	err = db.Get(&invoice, "SELECT * FROM invoice WHERE number=?", number)
	if err!= nil {
		return nil ,err
	}

	rows, err :=d.GetInvoiceRowsByInvoiceId(invoice.Id)
	if err!=nil {
		return nil, err
	}
	invoice.rows = rows

	return &invoice, nil
}

func (d *database) GetAllCustomers() ([]Customer,error) {
	db, err := d.getDatabase()
	if err!= nil {
		return nil ,err
	}
	customers := []Customer{}
	err = db.Select(&customers, "SELECT * FROM customer")
	if err!=nil {
		return nil, err
	}

	return customers, nil
}

func (d *database) GetCustomerByNumber(number string) (*Customer,error) {
	db, err := d.getDatabase()
	if err!= nil {
		return nil ,err
	}
	customer := Customer{}
	err = db.Get(&customer, "SELECT * FROM customer WHERE number=?", number)
	if err!= nil {
		return nil ,err
	}
	return &customer, nil
}

func (d *database) GetAllProducts() ([]Product,error) {
	db, err := d.getDatabase()
	if err!= nil {
		return nil ,err
	}
	products := []Product{}
	err = db.Select(&products, "SELECT * FROM product")
	if err!=nil {
		return nil, err
	}

	return products, nil
}

func (d *database) GetProductByNumber(number string) (*Product,error) {
	db, err := d.getDatabase()
	if err!= nil {
		return nil ,err
	}
	product := Product{}
	err = db.Get(&product, "SELECT * FROM product WHERE number=?", number)
	if err!= nil {
		return nil ,err
	}
	return &product, nil
}

func (d *database) getDatabase() (*sqlx.DB, error) {
	if d.db==nil {
		db, err := sqlx.Connect("mysql", "per:KnaskimGjwQ6M!@tcp(192.168.1.3:3306)/softinvoice?parseTime=True")
		if err != nil {
			return nil, err
		}
		d.db = db
	}
	return d.db, nil
}

func (d *database) CloseDatabase() {
	if d.db==nil {
		return
	}
	d.db.Close()
	d.db = nil
	return
}