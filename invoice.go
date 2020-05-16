package main

import (
	"database/sql"
	"time"
)

type Invoice struct {
	Id                    int
	Number                int
	Date                  time.Time
	DueDate               time.Time
	Amount                float32
	CustomerId            int
	CustomerNumber        string
	CustomerName          string
	CustomerAddress       string
	CustomerPostalAddress string
	CustomerReference     string
	PayDay                int
	Credit                bool
	CreditInvoiceNumber   sql.NullInt32
}
