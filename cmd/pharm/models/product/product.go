package product

import (
	"github.com/jinzhu/gorm"
	"github.com/warete/pharm/cmd/pharm/database"
)

type Product struct {
	gorm.Model
	Guid   string `json:"guid"`
	Name   string `json:"name"`
	Vendor string `json:"vendor"`
	ATH    string `json:"ath"`
	MNN    string `json:"mnn"`
	Active bool   `json:"active"`
}

func GetAll() ([]Product, error) {
	var products []Product
	database.DB.Connection.Find(&products)

	return products, nil
}

func GetById(id int) (Product, error) {
	var product Product
	database.DB.Connection.Find(&product, id)

	return product, nil
}

func Add(prod *Product) {
	database.DB.Connection.Create(&prod)
}

func Update(prod *Product) {
	database.DB.Connection.Save(&prod)
}
