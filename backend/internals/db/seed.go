package db

import (
	"aetrix/observer/internals/models"
	"log"

	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	// clear existing data
	db.Exec("DELETE FROM machine_users") // junction table if exists
	db.Exec("DELETE FROM machines")
	db.Exec("DELETE FROM users")

	// seed user
	user := models.User{Username: "user1", Password: "user1pass"}
	if err := db.Create(&user).Error; err != nil {
		panic("failed to seed users")
	}

	// seed machines
	machines := []models.Machine{
		{CreatorID: user.ID, Name: "machine1", IP: "192.168.1.1", Users: []models.User{user}},
		{CreatorID: user.ID, Name: "machine2", IP: "192.168.1.2", Users: []models.User{user}},
		{CreatorID: user.ID, Name: "machine3", IP: "192.168.1.3", Users: []models.User{user}},
	}

	for i := range machines {
		if err := db.Create(&machines[i]).Error; err != nil {
			panic("failed to seed machines")
		}
	}

	log.Println("Seeded Database")
}
