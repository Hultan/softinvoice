package database

import (
	"database/sql"
	"time"
)

type Invoice struct {
	Id                    int           `gorm:"column:id;primary_key"`
	Number                int           `gorm:"column:number"`
	Date                  time.Time     `gorm:"column:date"`
	DueDate               time.Time     `gorm:"column:duedate"`
	Amount                float32       `gorm:"column:amount"`
	CustomerNumber        string        `gorm:"column:customernumber;size:100"`
	CustomerName          string        `gorm:"column:customername;size:100"`
	CustomerAddress       string        `gorm:"column:customeraddress;size:100"`
	CustomerPostalAddress string        `gorm:"column:customerpostaladdress;size:100"`
	CustomerReference     string        `gorm:"column:customerreference;size:100"`
	PayDay                int           `gorm:"column:payday"`
	Credit                bool          `gorm:"column:credit;default:0"`
	CreditInvoiceNumber   sql.NullInt32 `gorm:"column:creditinvoicenumber;default:null"`
	ReadOnly              bool          `gorm:"column:readonly;default:false"`

	Rows []InvoiceRow
}

func (i *Invoice) TableName() string {
	return "invoice"
}
