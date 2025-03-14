package main

import (
	"context"
	"log"

	"redditBack/handler"
	"redditBack/model"
	"redditBack/repository"
	"redditBack/service"
	"redditBack/utility"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	db := connetToPostgreSQL()
	rdb := connetToRedis()

	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	voteRepo := repository.NewVoteRepository(db)
	cacheRepo := repository.NewRedisCacheRepository(rdb)

	authService := service.NewAuthService(&userRepo, &cacheRepo)
	postService := service.NewPostService(&postRepo, &userRepo)
	voteService := service.NewVoteService(&voteRepo, &postRepo, &userRepo, &cacheRepo)

	util := utility.NewUtility(&cacheRepo)

	authHandler := handler.NewAuthHandler(authService)
	postHandler := handler.NewPostHandler(postService)
	voteHandler := handler.NewVoteHandler(voteService)

	router := gin.Default()
	router.POST("/signup", authHandler.SignUp)
	router.POST("/login", authHandler.Login)
	auth := router.Group("/")
	auth.Use(util.JWTAuthMiddleware())
	{
		auth.POST("/signout", authHandler.SignOut)
		auth.POST("/posts/create", postHandler.CreatePost)
		auth.PUT("/posts/update", postHandler.EditPost)
		auth.DELETE("/posts/remove", postHandler.RemovePost)
		auth.POST("/vote", voteHandler.VotePost)
	}
	router.Run("0.0.0.0:8080")
}

func connetToPostgreSQL() *gorm.DB {
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
	return db
}

func connetToRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{

		Addr: "0.0.0.0:6380",

		Password: "",

		DB: 0,
	})
	status, err := rdb.Ping(context.Background()).Result()

	if err != nil {
		panic("Redis connection was refused")
	}
	log.Print(status)
	return rdb
}
