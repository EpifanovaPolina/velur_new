package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"online_store/database"
	"online_store/models"

	"github.com/gorilla/mux"
)

var (
	orderTpl        = template.Must(template.ParseFiles("templates/order.html"))
	orderSuccessTpl = template.Must(template.ParseFiles("templates/order_success.html"))
	adminOrdersTpl  = template.Must(template.ParseFiles("templates/order_list.html"))
)

type OrderHandler struct {
	Storage      *database.Storage
	ProductStore *models.ProductStore
}

func RegisterOrderRoutes(r *mux.Router, storage *database.Storage, productStore *models.ProductStore) {
	h := &OrderHandler{Storage: storage, ProductStore: productStore}
	r.HandleFunc("/order/{id:[0-9]+}", h.OrderHandler).Methods("GET")
	r.HandleFunc("/order", h.SubmitOrderHandler).Methods("POST")
	r.HandleFunc("/order_success", h.OrderSuccessHandler).Methods("GET")
	r.HandleFunc("/admin/orders", h.AdminOrdersHandler).Methods("GET")
	r.HandleFunc("/admin/order/{id}/delete", h.DeleteOrderHandler).Methods("POST")
}

func (h *OrderHandler) OrderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	if product, exists := h.ProductStore.Clothes[productID]; exists {
		orderTpl.Execute(w, map[string]interface{}{
			"ProductName":     product.Name,
			"ProductID":       productID,
			"ProductCategory": "Одежда",
		})
		return
	}

	if product, exists := h.ProductStore.Accessories[productID]; exists {
		orderTpl.Execute(w, map[string]interface{}{
			"ProductName":     product.Name,
			"ProductID":       productID,
			"ProductCategory": "Аксессуар",
		})
		return
	}

	http.NotFound(w, r)
}

func (h *OrderHandler) SubmitOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	r.ParseForm()

	productName := r.FormValue("product_name")
	productCategory := r.FormValue("product_category")
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	middleName := r.FormValue("middle_name")
	phone := r.FormValue("phone")
	quantity, _ := strconv.Atoi(r.FormValue("quantity"))
	region := r.FormValue("region")
	city := r.FormValue("city")
	street := r.FormValue("street")
	house := r.FormValue("house")
	apartment := r.FormValue("apartment")

	if quantity <= 0 {
		quantity = 1
	}

	_, err := h.Storage.DB.Exec(`
        INSERT INTO orders 
        (product_name, product_category, first_name, last_name, middle_name, 
         phone, quantity, region, city, street, house, apartment)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		productName, productCategory, firstName, lastName, middleName,
		phone, quantity, region, city, street, house, apartment)

	if err != nil {
		log.Println("Ошибка сохранения заказа:", err)
		http.Error(w, "Ошибка", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/order_success", http.StatusSeeOther)
}

func (h *OrderHandler) OrderSuccessHandler(w http.ResponseWriter, r *http.Request) {
	orderSuccessTpl.Execute(w, nil)
}

func (h *OrderHandler) AdminOrdersHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := h.Storage.DB.Query(`
        SELECT id, product_name, product_category, first_name, last_name, middle_name,
               phone, quantity, region, city, street, house, apartment, order_date
        FROM orders ORDER BY order_date DESC`)
	if err != nil {
		log.Println("Ошибка загрузки заказов:", err)
		http.Error(w, "Ошибка", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type OrderItem struct {
		ID, Quantity int
		ProductName, ProductCategory, FirstName, LastName, MiddleName, Phone,
		Region, City, Street, House, Apartment, OrderDate string
	}

	var orders []OrderItem
	for rows.Next() {
		var o OrderItem
		rows.Scan(&o.ID, &o.ProductName, &o.ProductCategory, &o.FirstName, &o.LastName,
			&o.MiddleName, &o.Phone, &o.Quantity, &o.Region, &o.City, &o.Street,
			&o.House, &o.Apartment, &o.OrderDate)
		orders = append(orders, o)
	}
	adminOrdersTpl.Execute(w, orders)
}

func (h *OrderHandler) DeleteOrderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	h.Storage.DB.Exec("DELETE FROM orders WHERE id=$1", id)
	http.Redirect(w, r, "/admin/orders", http.StatusSeeOther)
}
