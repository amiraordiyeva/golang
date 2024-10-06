package main

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID   uint `gorm:"primaryKey"`
	Name string
	Age  int
}

func insertUsers(db *gorm.DB) {
	user1 := User{Name: "Amira", Age: 19}
	user2 := User{Name: "Kamila", Age: 20}
	db.Create(&user1)
	db.Create(&user2)

	log.Println("Users inserted")
}

func retrieveUsers(db *gorm.DB) {
	var users []User
	db.Find(&users)

	for _, user := range users {
		log.Printf("ID: %d, Name: %s, Age: %d\n", user.ID, user.Name, user.Age)
	}
}

func main() {
	connStr := "user=amiraordiyeva dbname=golang sslmode=disable"

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal("Failed to migrate the database:", err)
	}

	insertUsers(db)
	retrieveUsers(db)
}
