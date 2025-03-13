package main

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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

type Post struct {
	ID          uint      `gorm:"primaryKey"`
	Title       string    `gorm:"not null"`
	Content     string    `gorm:"not null;type:text"`
	UserID      uint      `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	CachedScore int       `gorm:"default:0"`
	User        User      `gorm:"foreignKey:UserID"`
	Votes       []Vote    `gorm:"foreignKey:PostID"`
}

type Vote struct {
	UserID    uint      `gorm:"primaryKey"`
	PostID    uint      `gorm:"primaryKey"`
	VoteValue int       `gorm:"check:vote_value IN (-1,1)"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Post      Post      `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
}

func main() {
	dsn := "host=localhost user=pg password=pass dbname=reddit port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	err = db.AutoMigrate(&User{}, &Post{}, &Vote{})
	if err != nil {
		panic("Migration failed")
	}

	migrator := db.Migrator()
    
    tables := []string{"users", "posts", "votes"}
    for _, table := range tables {
        exists := migrator.HasTable(table)
        if exists {
            log.Printf("Table %s exists", table)
        } else {
            log.Printf("Table %s does NOT exist", table)
			panic("Error in creating tables , exiting...")
        }
    }

	log.Print("database connection successful , created tables")
}
