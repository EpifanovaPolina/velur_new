package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"online_store/database"
	"online_store/models"

	"github.com/gorilla/mux"
)

var (
	marketTpl        = template.Must(template.ParseFiles("templates/market.html"))
	aboutTpl         = template.Must(template.ParseFiles("templates/about.html"))
	productTpl       = template.Must(template.ParseFiles("templates/product.html"))
	accessoryTpl     = template.Must(template.ParseFiles("templates/accessory.html"))
	adminTpl         = template.Must(template.ParseFiles("templates/admin.html"))
	addProductTpl    = template.Must(template.ParseFiles("templates/add_product.html"))
	editClothingTpl  = template.Must(template.ParseFiles("templates/edit_clothing.html"))
	editAccessoryTpl = template.Must(template.ParseFiles("templates/edit_accessory.html"))
)

type ProductHandler struct {
	Storage      *database.Storage
	ProductStore *models.ProductStore
}

func RegisterProductRoutes(r *mux.Router, storage *database.Storage, productStore *models.ProductStore) {
	h := &ProductHandler{Storage: storage, ProductStore: productStore}
	r.HandleFunc("/", h.MarketHandler).Methods("GET")
	r.HandleFunc("/about", h.AboutHandler).Methods("GET")
	r.HandleFunc("/admin", h.AdminHandler).Methods("GET")
	r.HandleFunc("/clothing/{id:[0-9]+}", h.ClothingHandler).Methods("GET")
	r.HandleFunc("/accessory/{id:[0-9]+}", h.AccessoryHandler).Methods("GET")
	r.HandleFunc("/admin/add-clothing", h.AddClothingHandler).Methods("POST")
	r.HandleFunc("/admin/add-accessory", h.AddAccessoryHandler).Methods("POST")
	r.HandleFunc("/admin/delete-clothing/{id:[0-9]+}", h.DeleteClothingHandler).Methods("POST")
	r.HandleFunc("/admin/delete-accessory/{id:[0-9]+}", h.DeleteAccessoryHandler).Methods("POST")
	r.HandleFunc("/edit/clothing/{id:[0-9]+}", h.EditClothingForm).Methods("GET")
	r.HandleFunc("/edit/clothing/{id:[0-9]+}", h.EditClothingSubmit).Methods("POST")
	r.HandleFunc("/edit/accessory/{id:[0-9]+}", h.EditAccessoryForm).Methods("GET")
	r.HandleFunc("/edit/accessory/{id:[0-9]+}", h.EditAccessorySubmit).Methods("POST")
}

func (h *ProductHandler) MarketHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Clothes":     h.ProductStore.Clothes,
		"Accessories": h.ProductStore.Accessories,
		"Username":    "",
		"Role":        "",
	}
	marketTpl.Execute(w, data)
}

func (h *ProductHandler) AboutHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Username": "",
		"Role":     "",
	}
	aboutTpl.Execute(w, data)
}

func (h *ProductHandler) AdminHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Clothes":     h.ProductStore.Clothes,
		"Accessories": h.ProductStore.Accessories,
		"Username":    "",
		"Role":        "admin",
	}
	adminTpl.Execute(w, data)
}

func (h *ProductHandler) ClothingHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	product, exists := h.ProductStore.Clothes[id]
	if !exists {
		http.NotFound(w, r)
		return
	}
	data := map[string]interface{}{
		"Product":  product,
		"Username": "",
		"Role":     "",
	}
	productTpl.Execute(w, data)
}

func (h *ProductHandler) AccessoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	product, exists := h.ProductStore.Accessories[id]
	if !exists {
		http.NotFound(w, r)
		return
	}
	data := map[string]interface{}{
		"Product":  product,
		"Username": "",
		"Role":     "",
	}
	accessoryTpl.Execute(w, data)
}

func (h *ProductHandler) AddClothingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	r.ParseForm()
	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
	var id int
	h.Storage.DB.QueryRow(`
        INSERT INTO clothes (name, description, price, size, color, material, type, season)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		r.FormValue("name"), r.FormValue("description"), price,
		r.FormValue("size"), r.FormValue("color"), r.FormValue("material"),
		r.FormValue("type"), r.FormValue("season")).Scan(&id)

	h.ProductStore.Clothes[fmt.Sprint(id)] = models.Clothing{
		ID: fmt.Sprint(id), Name: r.FormValue("name"), Description: r.FormValue("description"),
		Price: price, Size: r.FormValue("size"), Color: r.FormValue("color"),
		Material: r.FormValue("material"), Type: r.FormValue("type"), Season: r.FormValue("season"),
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (h *ProductHandler) AddAccessoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	r.ParseForm()
	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
	var id int
	h.Storage.DB.QueryRow(`
        INSERT INTO accessories (name, description, price, type, color, material)
        VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		r.FormValue("name"), r.FormValue("description"), price,
		r.FormValue("type"), r.FormValue("color"), r.FormValue("material")).Scan(&id)

	h.ProductStore.Accessories[fmt.Sprint(id)] = models.Accessory{
		ID: fmt.Sprint(id), Name: r.FormValue("name"), Description: r.FormValue("description"),
		Price: price, Type: r.FormValue("type"), Color: r.FormValue("color"),
		Material: r.FormValue("material"),
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (h *ProductHandler) DeleteClothingHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	h.Storage.DB.Exec("DELETE FROM clothes WHERE id=$1", id)
	delete(h.ProductStore.Clothes, id)
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (h *ProductHandler) DeleteAccessoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	h.Storage.DB.Exec("DELETE FROM accessories WHERE id=$1", id)
	delete(h.ProductStore.Accessories, id)
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (h *ProductHandler) EditClothingForm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	product, exists := h.ProductStore.Clothes[id]
	if !exists {
		http.NotFound(w, r)
		return
	}
	editClothingTpl.Execute(w, product)
}

func (h *ProductHandler) EditClothingSubmit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	r.ParseForm()
	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
	h.Storage.DB.Exec(`
        UPDATE clothes SET name=$1, description=$2, price=$3, size=$4, 
        color=$5, material=$6, type=$7, season=$8 WHERE id=$9`,
		r.FormValue("name"), r.FormValue("description"), price,
		r.FormValue("size"), r.FormValue("color"), r.FormValue("material"),
		r.FormValue("type"), r.FormValue("season"), id)

	old := h.ProductStore.Clothes[id]
	h.ProductStore.Clothes[id] = models.Clothing{
		ID: id, Name: r.FormValue("name"), Description: r.FormValue("description"),
		Price: price, ImageURL: old.ImageURL, Size: r.FormValue("size"),
		Color: r.FormValue("color"), Material: r.FormValue("material"),
		Type: r.FormValue("type"), Season: r.FormValue("season"),
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (h *ProductHandler) EditAccessoryForm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	product, exists := h.ProductStore.Accessories[id]
	if !exists {
		http.NotFound(w, r)
		return
	}
	editAccessoryTpl.Execute(w, product)
}

func (h *ProductHandler) EditAccessorySubmit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	r.ParseForm()
	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
	h.Storage.DB.Exec(`
        UPDATE accessories SET name=$1, description=$2, price=$3, type=$4, color=$5, material=$6 WHERE id=$7`,
		r.FormValue("name"), r.FormValue("description"), price,
		r.FormValue("type"), r.FormValue("color"), r.FormValue("material"), id)

	old := h.ProductStore.Accessories[id]
	h.ProductStore.Accessories[id] = models.Accessory{
		ID: id, Name: r.FormValue("name"), Description: r.FormValue("description"),
		Price: price, ImageURL: old.ImageURL, Type: r.FormValue("type"),
		Color: r.FormValue("color"), Material: r.FormValue("material"),
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
