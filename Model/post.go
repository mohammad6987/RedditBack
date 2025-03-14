package model

import "time"

type Post struct {
	ID          uint      `gorm:"primaryKey"`
	Title       string    `gorm:"not null"`
	Content     string    `gorm:"not null;type:text"`
	UserID      uint      `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	CachedScore int       `gorm:"default:0"`
	User        User      `gorm:"foreignKey:UserID"`
	Votes       []Vote    `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
}
