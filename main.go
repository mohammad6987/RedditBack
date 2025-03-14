package main

import (
	"context"
	"log"

	"redditBack/handler"
	"redditBack/model"
	"redditBack/repository"
	"redditBack/service"
	"redditBack/utility"

	_ "redditBack/docs"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title           Reddit Clone API
// @version         1.0
// @description     API documentation for Reddit-like application
// @termsOfService  http://swagger.io/terms/
// @contact.name    API Support
// @contact.url     http://www.example.com/support
// @contact.email   support@example.com
// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html
// @host            localhost:8080
// @BasePath        /api/v1
// @in              header
// @name            Authorization
// @schemes         http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @tag.name posts
// @tag.description Post management operations
// @tag.name votes
// @tag.description Post voting operations
func main() {

	db := connetToPostgreSQL()
	rdb := connetToRedis()

	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	voteRepo := repository.NewVoteRepository(db)
	cacheRepo := repository.NewRedisCacheRepository(rdb)

	authService := service.NewAuthService(&userRepo, &cacheRepo)
	postService := service.NewPostService(&postRepo, &userRepo, &cacheRepo, &voteRepo)
	voteService := service.NewVoteService(&voteRepo, &postRepo, &userRepo, &cacheRepo)

	util := utility.NewUtility(&cacheRepo)

	authHandler := handler.NewAuthHandler(authService)
	postHandler := handler.NewPostHandler(postService)
	voteHandler := handler.NewVoteHandler(voteService)

	router := gin.Default()
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.POST("/signup", authHandler.SignUp)
	router.POST("/login", authHandler.Login)
	auth := router.Group("/")
	auth.Use(util.JWTAuthMiddleware())
	{
		auth.GET("/top", postHandler.GetTopPosts)
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
