package database

import (
	"log"

	"gorm.io/gorm"
)

// Erstellt eine neue User-Instanz
func CreateUser(db *gorm.DB, user *User) error {
	if err := db.Create(user).Error; err != nil {
		log.Printf("failed to create user: %v", err)
		return err
	}
	log.Println("User created successfully.")
	return nil
}

// Löscht einen User nach ID
func DeleteUserByID(db *gorm.DB, id uint) error {
	if err := db.Delete(&User{}, id).Error; err != nil {
		log.Printf("failed to delete user: %v", err)
		return err
	}
	log.Println("User deleted successfully.")
	return nil
}

// Findet alle User
func FindAllUsers(db *gorm.DB) ([]User, error) {
	var users []User
	if err := db.Find(&users).Error; err != nil {
		log.Printf("failed to find all users: %v", err)
		return nil, err
	}
	log.Println("Found all users successfully.")
	return users, nil
}
