package database

type InvoiceRow struct {
	Id        int     `gorm:"column:id;primary_key"`
	InvoiceId int     `gorm:"column:invoiceid;foreignkey:id"`
	Text      string  `gorm:"column:producttext;size:100"`
	Name      string  `gorm:"column:productname;size:100"`
	Price     float32 `gorm:"column:productprice"`
	Amount    float32 `gorm:"column:amount"`
	Total     float32 `gorm:"column:rowtotal"`
}

func (p *InvoiceRow) TableName() string {
	return "invoicerow"
}
