package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"redditBack/model"
)

func main() {
	dsn := "host=localhost user=pg password=pass dbname=reddit port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	err = db.AutoMigrate(&model.User{}, &model.Post{}, &model.Vote{})
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

	rdb := redis.NewClient(&redis.Options{

		Addr: "0.0.0.0:6380",

		Password: "",

		DB: 0,
	})
	defer rdb.Close()
	status, err := rdb.Ping(context.Background()).Result()

	if err != nil {
		panic("Redis connection was refused")
	}
	log.Print(status)

	router := gin.Default()
	router.POST("/signup", signUp)
	router.Run("0.0.0.0:8080")
}

func signUp(c *gin.Context) {
	type SignupRequest struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

}
