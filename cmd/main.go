package main

import (
	"os"

	"seojoonrp/board-api/internal/database"
	"seojoonrp/board-api/internal/handler"
	"seojoonrp/board-api/internal/repository"
	"seojoonrp/board-api/internal/server"
	"seojoonrp/board-api/internal/service"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// TODO: Separate config logic
	_ = godotenv.Load()
	uri := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB")
	jwtSecret := os.Getenv("JWT_SECRET")

	db, err := database.Connect(uri, dbName)
	if err != nil {
		panic(err)
	}
	defer database.Disconnect(db)

	postRepo := repository.NewPostRepo(db.Collection("posts"))
	userRepo := repository.NewUserRepo(db.Collection("users"))

	postSvc := service.NewPostService(postRepo, userRepo)
	userSvc := service.NewUserService(userRepo, jwtSecret)

	postHandler := handler.NewPostHandler(postSvc)
	userHandler := handler.NewUserHandler(userSvc)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	server.SetupRoutes(e, postHandler, userHandler, jwtSecret)

	e.Logger.Fatal(e.Start(":8080"))
}
