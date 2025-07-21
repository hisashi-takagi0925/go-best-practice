package database

import (
	"log"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	log.Println("Running database migrations...")

	err := AutoMigrate(db)
	if err != nil {
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}

func SeedData(db *gorm.DB) error {
	log.Println("Seeding initial data...")

	var count int64
	db.Model(&User{}).Count(&count)
	if count > 0 {
		log.Println("Data already exists, skipping seed")
		return nil
	}

	users := []User{
		{
			Name:     "John Doe",
			Username: "johndoe",
			Email:    "john@example.com",
		},
		{
			Name:     "Jane Smith",
			Username: "janesmith",
			Email:    "jane@example.com",
		},
	}

	for _, user := range users {
		if err := db.Create(&user).Error; err != nil {
			return err
		}
	}

	posts := []Post{
		{
			UserID: 1,
			Title:  "First Post",
			Body:   "This is the content of the first post.",
		},
		{
			UserID: 1,
			Title:  "Second Post",
			Body:   "This is the content of the second post.",
		},
		{
			UserID: 2,
			Title:  "Jane's Post",
			Body:   "This is Jane's first post.",
		},
	}

	for _, post := range posts {
		if err := db.Create(&post).Error; err != nil {
			return err
		}
	}

	log.Println("Initial data seeded successfully")
	return nil
}