package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func connectToDB() (*sql.DB, error) {
	connStr := "user=amiraordiyeva dbname=golang sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func createUserTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL,
        age INT
    )`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Table created")
}

func insertUser(db *sql.DB, name string, age int) {
	query := `INSERT INTO users (name, age) VALUES ($1, $2)`
	_, err := db.Exec(query, name, age)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("User inserted")
}

func queryAllUsers(db *sql.DB) {
	rows, err := db.Query("SELECT id, name, age FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var age int
		err := rows.Scan(&id, &name, &age)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %d, Name: %s, Age: %d\n", id, name, age)
	}
}

func main() {
	db, err := connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createUserTable(db)
	insertUser(db, "Amira", 19)
	insertUser(db, "Kamila", 20)
	queryAllUsers(db)
}
