package main

import (
	"os"

	"seojoonrp/board-api/internal/database"
	"seojoonrp/board-api/internal/handler"
	"seojoonrp/board-api/internal/repository"
	"seojoonrp/board-api/internal/server"
	"seojoonrp/board-api/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	uri := getenv("MONGO_URI", "mongodb://localhost:27017")
	dbName := getenv("MONGO_DB", "board")

	db, err := database.Connect(uri, dbName)
	if err != nil {
		panic(err)
	}
	defer database.Disconnect(db)

	postRepo := repository.NewPostRepo(db.Collection("posts"))
	userRepo := repository.NewUserRepo(db.Collection("users"))

	postSvc := service.NewPostService(postRepo, userRepo)
	userSvc := service.NewUserService(userRepo)

	postHandler := handler.NewPostHandler(postSvc)
	userHandler := handler.NewUserHandler(userSvc)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	jwtSecret := getenv("JWT_SECRET", "dev-jwt-secret")
	server.SetupRoutes(e, postHandler, userHandler, jwtSecret)

	e.Logger.Fatal(e.Start(":8080"))
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
