package model

import "time"

type Vote struct {
	UserID    uint      `gorm:"primaryKey"`
	PostID    uint      `gorm:"primaryKey"`
	VoteValue int       `gorm:"check:vote_value IN (-1,1)"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Post      Post      `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
}
