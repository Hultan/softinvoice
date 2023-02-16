package database

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"github.com/hultan/softteam/framework"
)

type Database struct {
	UseTestDatabase bool
	db              *gorm.DB
}

func NewDatabase(useTestDatabase bool) *Database {
	database := new(Database)
	database.UseTestDatabase = useTestDatabase
	return database
}

func (d *Database) InsertInvoice(invoice *Invoice) error {
	db, err := d.getDatabase()
	if err != nil {
		return err
	}
	if result := db.Create(invoice); result.Error != nil {
		return result.Error
	}

	return nil
}

func (d *Database) UpdateProduct(product *Product) error {
	db, err := d.getDatabase()
	if err != nil {
		return err
	}
	if result := db.Save(product); result.Error != nil {
		return result.Error
	}

	return nil
}

func (d *Database) GetAllInvoices() ([]Invoice, error) {
	db, err := d.getDatabase()
	if err != nil {
		return nil, err
	}
	var invoices []Invoice
	if result := db.Order("number desc").Find(&invoices); result.Error != nil {
		return nil, result.Error
	}

	for index, invoice := range invoices {
		rows, err := d.GetInvoiceRowsByInvoiceId(invoice.Id)
		if err != nil {
			return nil, err
		}
		invoices[index].Rows = rows
	}

	return invoices, nil
}

func (d *Database) GetInvoiceRowsByInvoiceId(id int) ([]InvoiceRow, error) {
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

func (d *Database) GetInvoiceByNumber(number string) (*Invoice, error) {
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
	invoice.Rows = rows

	return &invoice, nil
}

func (d *Database) GetAllCustomers() ([]Customer, error) {
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

func (d *Database) GetCustomerByNumber(number string) (*Customer, error) {
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

func (d *Database) GetAllProducts() ([]Product, error) {
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

func (d *Database) GetProductByNumber(number string) (*Product, error) {
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

func (d *Database) GetNextInvoiceNumber() (int, error) {
	db, err := d.getDatabase()
	if err != nil {
		return -1, err
	}
	var result int

	row := db.Table("invoice").Select("MAX(number)").Row()
	err = row.Scan(&result)
	if err != nil {
		return -1, err
	}

	return result + 1, nil
}

func (d *Database) getDatabase() (*gorm.DB, error) {
	databaseName := d.GetDatabaseName()
	if d.db == nil {
		fw := framework.NewFramework()
		server, err := fw.IO.ReadAllText(credentialsPath)
		if err != nil {
			return nil, err
		}
		var connectionString = fmt.Sprintf("%s/%s?parseTime=True", server, databaseName)
		db, err := gorm.Open("mysql", connectionString)
		if err != nil {
			return nil, err
		}
		d.db = db
	}
	return d.db, nil
}

func (d *Database) CloseDatabase() {
	if d.db == nil {
		return
	}
	_ = d.db.Close()
	d.db = nil
	return
}

func (d *Database) GetDatabaseName() string {
	if d.UseTestDatabase {
		return constDatabaseNameTest
	} else {
		return constDatabaseName
	}
}
