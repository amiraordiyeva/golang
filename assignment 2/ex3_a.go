package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

var db *sql.DB

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	var err error
	connStr := "user=amiraordiyeva dbname=golang sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/users", handleUsers)
	fmt.Println("Сервер запущен на порту 8080")
	http.ListenAndServe(":8080", nil)
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getUsers(w, r)
	case "POST":
		createUser(w, r)
	}
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	age := r.URL.Query().Get("age")
	sort := r.URL.Query().Get("sort")

	var query string
	if age != "" {
		query = "SELECT * FROM users WHERE age = $1"
	} else {
		query = "SELECT * FROM users"
	}

	if sort == "name" {
		query += " ORDER BY name"
	}

	rows, err := db.Query(query, age)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Age); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Проверка уникальности имени
	var exists bool
	err := db.QueryRow("SELECT exists(SELECT 1 FROM users WHERE name = $1)", user.Name).Scan(&exists)
	if err != nil || exists {
		http.Error(w, "User with this name already exists", http.StatusConflict)
		return
	}

	_, err = db.Exec("INSERT INTO users (name, age) VALUES ($1, $2)", user.Name, user.Age)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
