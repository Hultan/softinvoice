package database

type Customer struct {
	Id            int    `gorm:"column:id;primary_key"`
	FancyName     string `gorm:"column:fancyname;size:100;default:null"`
	Name          string `gorm:"column:name;size:100"`
	Number        string `gorm:"column:number;size:100"`
	Reference     string `gorm:"column:reference;size:100:default:null"`
	PayDay        int    `gorm:"column:payday;default:30"`
	Address       string `gorm:"column:address;size:100"`
	PostalAddress string `gorm:"column:postaladdress;size:100"`
}

func (p *Customer) TableName() string {
	return "customer"
}
