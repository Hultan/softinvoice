package main

import (
	"database/sql"
	"time"
)

type Invoice struct {
	Id                    int 			`db:"id"`
	Number                int			`db:"number"`
	Date                  time.Time		`db:"date"`
	DueDate               time.Time		`db:"duedate"`
	Amount                float32		`db:"amount"`
	CustomerId            int			`db:"customerid"`
	CustomerNumber        string		`db:"customernumber"`
	CustomerName          string		`db:"customername"`
	CustomerAddress       string		`db:"customeraddress"`
	CustomerPostalAddress string		`db:"customerpostaladdress"`
	CustomerReference     string		`db:"customerreference"`
	PayDay                int			`db:"payday"`
	Credit                bool			`db:"credit"`
	CreditInvoiceNumber   sql.NullInt32 `db:"creditinvoicenumber"`
}
