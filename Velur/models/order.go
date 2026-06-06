package models

type Order struct {
	ProductName     string
	ProductCategory string
	FirstName       string
	LastName        string
	MiddleName      string
	Phone           string
	Quantity        int
	Region          string
	City            string
	Street          string
	House           string
	Apartment       string
}

type AdminOrder struct {
	ID              int
	ProductName     string
	ProductCategory string
	FirstName       string
	LastName        string
	MiddleName      string
	Phone           string
	Quantity        int
	Region          string
	City            string
	Street          string
	House           string
	Apartment       string
	OrderDate       string
}
