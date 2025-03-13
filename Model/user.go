package model

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey"`
	Username     string    `gorm:"unique;not null"`
	Email        string    `gorm:"unique;not null"`
	PasswordHash string    `gorm:"not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
	Posts        []Post    `gorm:"foreignKey:UserID"`
	Votes        []Vote    `gorm:"foreignKey:UserID"`
}
