package models

import "time"

type User struct {
	ID            uint   `gorm:"primaryKey"`
	Username      string `gorm:"unique;not null"`
	Email         string
	Password      string    `gorm:"not null"`
	OwnedMachines []Machine `gorm:"foreignKey:CreatorID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Machines      []Machine `gorm:"many2many:machine_users;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	CreatedAt time.Time
}
