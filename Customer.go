package main

type Customer struct {
	Id            int		`db:"id"`
	FancyName     string	`db:"fancyname"`
	Name          string	`db:"name"`
	Number        string	`db:"number"`
	Reference     string	`db:"reference"`
	PayDay        int		`db:"payday"`
	Address       string	`db:"address"`
	PostalAddress string	`db:"postaladdress"`
}
