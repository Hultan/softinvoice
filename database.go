package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type database struct {
	db *sqlx.DB
}

func (d *database) GetAllInvoices() []Invoice {
	db, err := d.getDatabase()
	invoices := []Invoice{}
	err = db.Select(&invoices, "SELECT * FROM softinvoice.invoice")
	if err!=nil {
		fmt.Println(err.Error())
		panic(err)
	}

	return invoices
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