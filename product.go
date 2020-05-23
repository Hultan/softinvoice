package main

type Product struct {
	Id     int     `gorm:"column:id;primary_key"`
	Number string  `gorm:"column:number;size:50"`
	Text   string  `gorm:"column:text;size:100"`
	Name   string  `gorm:"column:name;size:100"`
	Price  float32 `gorm:"column:price"`
}

func (p *Product) TableName() string {
	return "product"
}
