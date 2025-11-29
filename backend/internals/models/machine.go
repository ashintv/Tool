package models

import "time"

type Machine struct {
	ID        uint      `gorm:"primaryKey"`
	MachineID string    `gorm:"unique;not null"`

	CreatorID uint
	Creator   User `gorm:"foreignKey:CreatorID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Users     []User `gorm:"many2many:machine_users;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	IP        string
	CreatedAt time.Time
}

