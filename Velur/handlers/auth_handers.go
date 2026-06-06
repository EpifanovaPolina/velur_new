package handlers

import (
	"html/template"
	"net/http"
	"strings"

	"online_store/utils"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var (
	loginTpl        = template.Must(template.ParseFiles("templates/login.html"))
	registrationTpl = template.Must(template.ParseFiles("templates/registration.html"))
)

func RegisterAuthRoutes(r *mux.Router) {
	r.HandleFunc("/login", LoginHandler).Methods("GET", "POST")
	r.HandleFunc("/registration", RegistrationHandler).Methods("GET", "POST")
	r.HandleFunc("/logout", LogoutHandler).Methods("GET")
}

func GetUsername(r *http.Request) string {
	session, _ := utils.Store.Get(r, "session-name")
	if username, ok := session.Values["username"].(string); ok {
		return username
	}
	return ""
}

func GetUserRole(r *http.Request) string {
	session, _ := utils.Store.Get(r, "session-name")
	if role, ok := session.Values["role"].(string); ok {
		return role
	}
	return ""
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var hashedPassword, role string
		err := utils.DB.QueryRow("SELECT password, role FROM users WHERE username=$1", username).Scan(&hashedPassword, &role)
		if err != nil || bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) != nil {
			loginTpl.Execute(w, map[string]interface{}{"ErrorMessage": "Неверный логин или пароль"})
			return
		}

		session, _ := utils.Store.Get(r, "session-name")
		session.Values["username"] = username
		session.Values["role"] = role
		session.Save(r, w)

		if role == "admin" {
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		return
	}
	loginTpl.Execute(w, nil)
}

func RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := strings.TrimSpace(r.FormValue("username"))
		password := r.FormValue("password")
		email := strings.TrimSpace(r.FormValue("email"))

		if username == "" || password == "" || email == "" {
			registrationTpl.Execute(w, map[string]interface{}{"ErrorMessage": "Заполните все поля"})
			return
		}

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		role := "user"
		if username == "admin" {
			role = "admin"
		}

		_, err := utils.DB.Exec("INSERT INTO users (username, password, email, role) VALUES ($1, $2, $3, $4)",
			username, hashedPassword, email, role)
		if err != nil {
			registrationTpl.Execute(w, map[string]interface{}{"ErrorMessage": "Ошибка регистрации"})
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	registrationTpl.Execute(w, nil)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := utils.Store.Get(r, "session-name")
	delete(session.Values, "username")
	delete(session.Values, "role")
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
