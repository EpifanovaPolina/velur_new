package models

type Clothing struct {
	ID          string
	Name        string
	Description string
	ImageURL    string
	Price       float64
	Size        string
	Color       string
	Material    string
	Type        string
	Season      string
}

type Accessory struct {
	ID          string
	Name        string
	Description string
	ImageURL    string
	Price       float64
	Type        string
	Color       string
	Material    string
}
