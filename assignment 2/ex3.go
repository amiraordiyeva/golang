package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var sqlDB *sql.DB

type User struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"unique;not null"`
	Age  int    `json:"age"`
}

func initDB() {
	var err error
	dsn := "user=amiraordiyeva dbname=golang sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if err := db.AutoMigrate(&User{}); err != nil {
		panic(err)
	}

	sqlDB, err = sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
}

func handleError(c *gin.Context, err error) {
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func main() {
	initDB()
	defer sqlDB.Close()

	router := gin.Default()

	router.GET("/users", getUsersSQL)
	router.POST("/users", createUserSQL)
	router.PUT("/users/:id", updateUserSQL)
	router.DELETE("/users/:id", deleteUserSQL)

	router.Run(":8080")
}

func getUsersSQL(c *gin.Context) {
	var users []User
	rows, err := sqlDB.Query("SELECT id, name, age FROM users")
	handleError(c, err)
	defer rows.Close()

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Age); err != nil {
			handleError(c, err)
			return
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}

func createUserSQL(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := sqlDB.Exec("INSERT INTO users (name, age) VALUES ($1, $2)", user.Name, user.Age); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, user)
}

func updateUserSQL(c *gin.Context) {
	id := c.Param("id")
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := sqlDB.Exec("UPDATE users SET name=$1, age=$2 WHERE id=$3", user.Name, user.Age, id)
	handleError(c, err)

	c.JSON(http.StatusOK, user)
}

func deleteUserSQL(c *gin.Context) {
	id := c.Param("id")
	_, err := sqlDB.Exec("DELETE FROM users WHERE id=$1", id)
	handleError(c, err)

	c.Status(http.StatusNoContent)
}
