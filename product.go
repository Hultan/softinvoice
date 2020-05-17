package main

type Product struct {
	Id     int     `db:"id"`
	Number string  `db:"number"`
	Text   string  `db:"text"`
	Name   string  `db:"name"`
	Price  float32 `db:"price"`
}
