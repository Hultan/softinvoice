package main

type InvoiceRow struct {
	Id        int     `db:"id"`
	InvoiceId int     `db:"invoiceid"`
	Text      string  `db:"producttext"`
	Name      string  `db:"productname"`
	Price     float32 `db:"productprice"`
	Amount    float32 `db:"amount"`
	Total     float32 `db:"rowtotal"`
}

