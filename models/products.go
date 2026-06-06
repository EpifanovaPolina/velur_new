package models

import (
	"database/sql"
	"log"
	"strings"
)

type ProductStore struct {
	Clothes     map[string]Clothing
	Accessories map[string]Accessory
}

func NewProductStore() *ProductStore {
	return &ProductStore{
		Clothes:     make(map[string]Clothing),
		Accessories: make(map[string]Accessory),
	}
}

func (s *ProductStore) LoadProducts(db *sql.DB) {
	rows, err := db.Query(`
        SELECT id, name, description, image_url, price, 
               COALESCE(size, ''), COALESCE(color, ''), 
               COALESCE(material, ''), COALESCE(type, ''), COALESCE(season, '') 
        FROM clothes`)
	if err != nil {
		log.Println("Ошибка загрузки одежды:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var c Clothing
		var img sql.NullString
		rows.Scan(&c.ID, &c.Name, &c.Description, &img, &c.Price,
			&c.Size, &c.Color, &c.Material, &c.Type, &c.Season)
		if img.Valid {
			c.ImageURL = strings.ReplaceAll(img.String, "\\", "/")
		}
		s.Clothes[c.ID] = c
	}

	rows2, err := db.Query(`
        SELECT id, name, description, image_url, price, 
               COALESCE(type, ''), COALESCE(color, ''), COALESCE(material, '') 
        FROM accessories`)
	if err != nil {
		log.Println("Ошибка загрузки аксессуаров:", err)
		return
	}
	defer rows2.Close()

	for rows2.Next() {
		var a Accessory
		var img sql.NullString
		rows2.Scan(&a.ID, &a.Name, &a.Description, &img, &a.Price,
			&a.Type, &a.Color, &a.Material)
		if img.Valid {
			a.ImageURL = strings.ReplaceAll(img.String, "\\", "/")
		}
		s.Accessories[a.ID] = a
	}

	log.Printf("Загружено одежды: %d, аксессуаров: %d", len(s.Clothes), len(s.Accessories))
}
