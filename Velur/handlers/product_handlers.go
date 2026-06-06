package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"online_store/models"
	"online_store/utils"

	"github.com/gorilla/mux"
)

var (
	marketTpl        = template.Must(template.ParseFiles("templates/market.html"))
	aboutTpl         = template.Must(template.ParseFiles("templates/about.html"))
	productTpl       = template.Must(template.ParseFiles("templates/product.html"))
	accessoryTpl     = template.Must(template.ParseFiles("templates/accessory.html"))
	adminTpl         = template.Must(template.ParseFiles("templates/admin.html"))
	editClothingTpl  = template.Must(template.ParseFiles("templates/edit_clothing.html"))
	editAccessoryTpl = template.Must(template.ParseFiles("templates/edit_accessory.html"))
)

func RegisterProductRoutes(r *mux.Router) {
	r.HandleFunc("/", MarketHandler).Methods("GET")
	r.HandleFunc("/about", AboutHandler).Methods("GET")
	r.HandleFunc("/admin", AdminHandler).Methods("GET")
	r.HandleFunc("/clothing/{id}", ClothingHandler).Methods("GET")
	r.HandleFunc("/accessory/{id}", AccessoryHandler).Methods("GET")
	r.HandleFunc("/admin/add-clothing", AddClothingHandler).Methods("POST")
	r.HandleFunc("/admin/add-accessory", AddAccessoryHandler).Methods("POST")
	r.HandleFunc("/admin/delete-clothing/{id}", DeleteClothingHandler).Methods("POST")
	r.HandleFunc("/admin/delete-accessory/{id}", DeleteAccessoryHandler).Methods("POST")
	r.HandleFunc("/edit/clothing/{id}", EditClothingForm).Methods("GET")
	r.HandleFunc("/edit/clothing/{id}", EditClothingSubmit).Methods("POST")
	r.HandleFunc("/edit/accessory/{id}", EditAccessoryForm).Methods("GET")
	r.HandleFunc("/edit/accessory/{id}", EditAccessorySubmit).Methods("POST")
}

func MarketHandler(w http.ResponseWriter, r *http.Request) {
	utils.LoadProducts()
	data := map[string]interface{}{
		"Clothes":     models.Clothes,
		"Accessories": models.Accessories,
		"Username":    GetUsername(r),
		"Role":        GetUserRole(r),
	}
	marketTpl.Execute(w, data)
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Username": GetUsername(r),
		"Role":     GetUserRole(r),
	}
	aboutTpl.Execute(w, data)
}

func ClothingHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	log.Printf("ClothingHandler: id=%s", id)

	product, exists := models.Clothes[id]
	if !exists {
		http.NotFound(w, r)
		return
	}

	data := map[string]interface{}{
		"Product":  product,
		"Username": GetUsername(r),
		"Role":     GetUserRole(r),
	}
	productTpl.Execute(w, data)
}

func AccessoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	product, exists := models.Accessories[id]
	if !exists {
		http.NotFound(w, r)
		return
	}
	data := map[string]interface{}{
		"Product":  product,
		"Username": GetUsername(r),
		"Role":     GetUserRole(r),
	}
	accessoryTpl.Execute(w, data)
}

func AdminHandler(w http.ResponseWriter, r *http.Request) {
	if GetUserRole(r) != "admin" {
		http.Error(w, "Доступ запрещён", http.StatusForbidden)
		return
	}
	utils.LoadProducts()
	data := map[string]interface{}{
		"Clothes":     models.Clothes,
		"Accessories": models.Accessories,
		"Username":    GetUsername(r),
		"Role":        GetUserRole(r),
	}
	adminTpl.Execute(w, data)
}

func AddClothingHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
	var id int
	utils.DB.QueryRow(`
        INSERT INTO clothes (name, description, price, size, color, material, type, season)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		r.FormValue("name"), r.FormValue("description"), price,
		r.FormValue("size"), r.FormValue("color"), r.FormValue("material"),
		r.FormValue("type"), r.FormValue("season")).Scan(&id)
	models.Clothes[fmt.Sprint(id)] = models.Clothing{
		ID: fmt.Sprint(id), Name: r.FormValue("name"), Description: r.FormValue("description"),
		Price: price, Size: r.FormValue("size"), Color: r.FormValue("color"),
		Material: r.FormValue("material"), Type: r.FormValue("type"), Season: r.FormValue("season"),
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func AddAccessoryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
	var id int
	utils.DB.QueryRow(`
        INSERT INTO accessories (name, description, price, type, color, material)
        VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		r.FormValue("name"), r.FormValue("description"), price,
		r.FormValue("type"), r.FormValue("color"), r.FormValue("material")).Scan(&id)
	models.Accessories[fmt.Sprint(id)] = models.Accessory{
		ID: fmt.Sprint(id), Name: r.FormValue("name"), Description: r.FormValue("description"),
		Price: price, Type: r.FormValue("type"), Color: r.FormValue("color"),
		Material: r.FormValue("material"),
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func DeleteClothingHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	utils.DB.Exec("DELETE FROM clothes WHERE id=$1", id)
	delete(models.Clothes, id)
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func DeleteAccessoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	utils.DB.Exec("DELETE FROM accessories WHERE id=$1", id)
	delete(models.Accessories, id)
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func EditClothingForm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	product, exists := models.Clothes[id]
	if !exists {
		http.NotFound(w, r)
		return
	}
	editClothingTpl.Execute(w, product)
}

func EditClothingSubmit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	r.ParseForm()
	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
	utils.DB.Exec(`
        UPDATE clothes SET name=$1, description=$2, price=$3, size=$4, 
        color=$5, material=$6, type=$7, season=$8 WHERE id=$9`,
		r.FormValue("name"), r.FormValue("description"), price,
		r.FormValue("size"), r.FormValue("color"), r.FormValue("material"),
		r.FormValue("type"), r.FormValue("season"), id)
	old := models.Clothes[id]
	models.Clothes[id] = models.Clothing{
		ID: id, Name: r.FormValue("name"), Description: r.FormValue("description"),
		Price: price, ImageURL: old.ImageURL, Size: r.FormValue("size"),
		Color: r.FormValue("color"), Material: r.FormValue("material"),
		Type: r.FormValue("type"), Season: r.FormValue("season"),
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func EditAccessoryForm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	product, exists := models.Accessories[id]
	if !exists {
		http.NotFound(w, r)
		return
	}
	editAccessoryTpl.Execute(w, product)
}

func EditAccessorySubmit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	r.ParseForm()
	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
	utils.DB.Exec(`
        UPDATE accessories SET name=$1, description=$2, price=$3, type=$4, color=$5, material=$6 WHERE id=$7`,
		r.FormValue("name"), r.FormValue("description"), price,
		r.FormValue("type"), r.FormValue("color"), r.FormValue("material"), id)
	old := models.Accessories[id]
	models.Accessories[id] = models.Accessory{
		ID: id, Name: r.FormValue("name"), Description: r.FormValue("description"),
		Price: price, ImageURL: old.ImageURL, Type: r.FormValue("type"),
		Color: r.FormValue("color"), Material: r.FormValue("material"),
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
