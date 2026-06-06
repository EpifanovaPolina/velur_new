package utils

import (
	"database/sql"
	"log"
	"strings"

	"online_store/models"
)

func LoadProducts() {
	rows, err := DB.Query(`
        SELECT id, name, description, image_url, price, 
               COALESCE(size, ''), COALESCE(color, ''), 
               COALESCE(material, ''), COALESCE(type, ''), COALESCE(season, '') 
        FROM clothes`)
	if err != nil {
		log.Println("Ошибка загрузки одежды:", err)
		return
	}
	defer rows.Close()

	models.Clothes = make(map[string]models.Clothing)
	for rows.Next() {
		var c models.Clothing
		var img sql.NullString
		rows.Scan(&c.ID, &c.Name, &c.Description, &img, &c.Price,
			&c.Size, &c.Color, &c.Material, &c.Type, &c.Season)
		if img.Valid {
			c.ImageURL = strings.ReplaceAll(img.String, "\\", "/")
		}
		models.Clothes[c.ID] = c
	}

	rows2, _ := DB.Query(`
        SELECT id, name, description, image_url, price, 
               COALESCE(type, ''), COALESCE(color, ''), COALESCE(material, '') 
        FROM accessories`)
	defer rows2.Close()
	models.Accessories = make(map[string]models.Accessory)
	for rows2.Next() {
		var a models.Accessory
		var img sql.NullString
		rows2.Scan(&a.ID, &a.Name, &a.Description, &img, &a.Price,
			&a.Type, &a.Color, &a.Material)
		if img.Valid {
			a.ImageURL = strings.ReplaceAll(img.String, "\\", "/")
		}
		models.Accessories[a.ID] = a
	}
	log.Printf("Загружено одежды: %d, аксессуаров: %d", len(models.Clothes), len(models.Accessories))
}
