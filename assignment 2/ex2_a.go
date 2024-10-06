package main

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID      uint   `gorm:"primaryKey"`
	Name    string `gorm:"not null"`
	Age     int    `gorm:"not null"`
	Profile Profile
}

type Profile struct {
	ID                uint `gorm:"primaryKey"`
	UserID            uint `gorm:"not null;unique"`
	Bio               string
	ProfilePictureURL string
}

func createUserWithProfile(db *gorm.DB, user User, profile Profile) uint {
	if err := db.Create(&user).Error; err != nil {
		log.Fatal("Ошибка вставки пользователя:", err)
	}
	profile.UserID = user.ID
	if err := db.Create(&profile).Error; err != nil {
		log.Fatal("Ошибка вставки профиля:", err)
	}
	log.Println("Пользователь и профиль успешно вставлены.")
	return user.ID // Возвращаем ID созданного пользователя
}

func queryUsersWithProfiles(db *gorm.DB) {
	var users []User
	if err := db.Preload("Profile").Find(&users).Error; err != nil {
		log.Fatal("Ошибка запроса пользователей с профилями:", err)
	}
	for _, user := range users {
		log.Printf("ID: %d, Name: %s, Age: %d, Bio: %s\n", user.ID, user.Name, user.Age, user.Profile.Bio)
	}
}

func updateUserProfile(db *gorm.DB, userID uint, newBio string) {
	var profile Profile
	if err := db.First(&profile, "user_id = ?", userID).Error; err != nil {
		log.Fatal("Ошибка получения профиля:", err)
	}
	profile.Bio = newBio
	if err := db.Save(&profile).Error; err != nil {
		log.Fatal("Ошибка обновления профиля:", err)
	}
	log.Println("Профиль обновлен.")
}

func deleteUserWithProfile(db *gorm.DB, userID uint) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var user User
	if err := tx.First(&user, userID).Error; err != nil {
		log.Fatal("Ошибка получения пользователя:", err)
	}

	var profile Profile
	if err := tx.First(&profile, "user_id = ?", userID).Error; err == nil {
		if err := tx.Delete(&profile).Error; err != nil {
			log.Fatal("Ошибка удаления профиля:", err)
		}
	}

	if err := tx.Delete(&user).Error; err != nil {
		log.Fatal("Ошибка удаления пользователя:", err)
	}

	tx.Commit()
	log.Println("Пользователь и профиль удалены.")
}

func main() {
	dsn := "user=amiraordiyeva dbname=golang sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}

	log.Println("Подключение к базе данных успешно!")

	err = db.AutoMigrate(&User{}, &Profile{})
	if err != nil {
		log.Fatal("Ошибка миграции базы данных:", err)
	}

	user := User{Name: "Amira", Age: 19}
	profile := Profile{Bio: "Программист", ProfilePictureURL: "url_to_picture"}
	userID := createUserWithProfile(db, user, profile)

	queryUsersWithProfiles(db)

	updateUserProfile(db, userID, "Обновленный био")

	deleteUserWithProfile(db, userID)
}
