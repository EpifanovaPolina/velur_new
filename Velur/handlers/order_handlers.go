package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"online_store/models"
	"online_store/utils"

	"github.com/gorilla/mux"
)

var (
	orderTpl        = template.Must(template.ParseFiles("templates/order.html"))
	orderSuccessTpl = template.Must(template.ParseFiles("templates/order_success.html"))
	adminOrdersTpl  = template.Must(template.ParseFiles("templates/order_list.html"))
)

func RegisterOrderRoutes(r *mux.Router) {
	r.HandleFunc("/order/{id}", OrderHandler).Methods("GET")
	r.HandleFunc("/order", SubmitOrderHandler).Methods("POST")
	r.HandleFunc("/order_success", OrderSuccessHandler).Methods("GET")
	r.HandleFunc("/admin/orders", AdminOrdersHandler).Methods("GET")
	r.HandleFunc("/admin/order/{id}/delete", DeleteOrderHandler).Methods("POST")
}

func OrderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	log.Printf("OrderHandler: ищу товар с ID=%s", productID)

	// Ищем в одежде
	if product, exists := models.Clothes[productID]; exists {
		log.Printf("Найдена одежда: %s", product.Name)
		data := map[string]interface{}{
			"ProductName":     product.Name,
			"ProductID":       productID,
			"ProductCategory": "Одежда",
		}
		orderTpl.Execute(w, data)
		return
	}

	// Ищем в аксессуарах
	if product, exists := models.Accessories[productID]; exists {
		log.Printf("Найден аксессуар: %s", product.Name)
		data := map[string]interface{}{
			"ProductName":     product.Name,
			"ProductID":       productID,
			"ProductCategory": "Аксессуар",
		}
		orderTpl.Execute(w, data)
		return
	}

	log.Printf("Товар с ID=%s не найден!", productID)
	http.NotFound(w, r)
}

func SubmitOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Ошибка обработки формы", http.StatusBadRequest)
		return
	}

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

	log.Printf("Сохранение заказа: %s, %s %s, тел: %s", productName, firstName, lastName, phone)

	_, err = utils.DB.Exec(`
        INSERT INTO orders 
        (product_name, product_category, first_name, last_name, middle_name, 
         phone, quantity, region, city, street, house, apartment)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		productName, productCategory, firstName, lastName, middleName,
		phone, quantity, region, city, street, house, apartment)

	if err != nil {
		log.Println("Ошибка при сохранении заказа:", err)
		http.Error(w, "Не удалось сохранить заказ", http.StatusInternalServerError)
		return
	}

	log.Println("Заказ успешно сохранён!")
	http.Redirect(w, r, "/order_success", http.StatusSeeOther)
}

func OrderSuccessHandler(w http.ResponseWriter, r *http.Request) {
	orderSuccessTpl.Execute(w, nil)
}

func AdminOrdersHandler(w http.ResponseWriter, r *http.Request) {
	if GetUserRole(r) != "admin" {
		http.Error(w, "Доступ запрещён", http.StatusForbidden)
		return
	}

	rows, err := utils.DB.Query(`
        SELECT id, product_name, product_category, first_name, last_name, middle_name,
               phone, quantity, region, city, street, house, apartment, order_date
        FROM orders ORDER BY order_date DESC`)
	if err != nil {
		log.Println("Ошибка загрузки заказов:", err)
		http.Error(w, "Ошибка загрузки заказов", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type OrderItem struct {
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

	var orders []OrderItem
	for rows.Next() {
		var o OrderItem
		err := rows.Scan(
			&o.ID, &o.ProductName, &o.ProductCategory,
			&o.FirstName, &o.LastName, &o.MiddleName,
			&o.Phone, &o.Quantity, &o.Region, &o.City,
			&o.Street, &o.House, &o.Apartment, &o.OrderDate,
		)
		if err != nil {
			log.Println("Ошибка сканирования:", err)
			continue
		}
		orders = append(orders, o)
	}

	adminOrdersTpl.Execute(w, orders)
}

func DeleteOrderHandler(w http.ResponseWriter, r *http.Request) {
	if GetUserRole(r) != "admin" {
		http.Error(w, "Доступ запрещён", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	_, err := utils.DB.Exec("DELETE FROM orders WHERE id = $1", id)
	if err != nil {
		log.Printf("Ошибка удаления заказа %s: %v", id, err)
		http.Error(w, "Ошибка при удалении заказа", http.StatusInternalServerError)
		return
	}

	log.Printf("Заказ %s удалён", id)
	http.Redirect(w, r, "/admin/orders", http.StatusSeeOther)
}
