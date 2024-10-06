package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type User struct {
	ID   int
	Name string
	Age  int
}

func main() {
	connStr := "user=amiraordiyeva dbname=golang sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Ошибка проверки подключения к базе данных:", err)
	}

	log.Println("Подключение к базе данных успешно!")

	createUsersTable(db)

	users := []User{
		{Name: "Amira", Age: 19},
		{Name: "Kamila", Age: 20},
	}
	insertMultipleUsers(db, users)

	queryUsers(db, 0, 1, 2)
	updateUserDetails(db, 1, "Amira Updated", 20)

	deleteUserByID(db, 2)
}

func createUsersTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) UNIQUE NOT NULL,
		age INT NOT NULL
	);`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Ошибка создания таблицы:", err)
	}
	log.Println("Таблица users создана.")
}

func insertMultipleUsers(db *sql.DB, users []User) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal("Ошибка начала транзакции:", err)
	}

	for _, user := range users {
		_, err := tx.Exec("INSERT INTO users (name, age) VALUES ($1, $2)", user.Name, user.Age)
		if err != nil {
			tx.Rollback()
			log.Fatal("Ошибка вставки пользователя:", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal("Ошибка подтверждения транзакции:", err)
	}
	log.Println("Пользователи вставлены.")
}

func queryUsers(db *sql.DB, ageFilter int, page int, pageSize int) {
	var query string
	var args []interface{}
	if ageFilter > 0 {
		query = "SELECT * FROM users WHERE age = $1 LIMIT $2 OFFSET $3"
		args = []interface{}{ageFilter, pageSize, (page - 1) * pageSize}
	} else {
		query = "SELECT * FROM users LIMIT $1 OFFSET $2"
		args = []interface{}{pageSize, (page - 1) * pageSize}
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		log.Fatal("Ошибка выполнения запроса:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Age); err != nil {
			log.Fatal("Ошибка сканирования строки:", err)
		}
		log.Printf("ID: %d, Name: %s, Age: %d\n", user.ID, user.Name, user.Age)
	}
}

func updateUserDetails(db *sql.DB, id int, name string, age int) {
	_, err := db.Exec("UPDATE users SET name = $1, age = $2 WHERE id = $3", name, age, id)
	if err != nil {
		log.Fatal("Ошибка обновления пользователя:", err)
	}
	log.Println("Пользователь обновлен.")
}

func deleteUserByID(db *sql.DB, id int) {
	_, err := db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		log.Fatal("Ошибка удаления пользователя:", err)
	}
	log.Println("Пользователь удален.")
}
