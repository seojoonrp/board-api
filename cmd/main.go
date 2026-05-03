package main

import (
	"log"
	"os"

	"seojoonrp/board-api/internal/apperror"
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
	e.HTTPErrorHandler = apperror.Handler
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogMethod:   true,
		LogLatency:  true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error != nil {
				log.Printf("%s %s %d %s ERROR: %v\n",
					v.Method, v.URI, v.Status, v.Latency, v.Error)
			} else {
				log.Printf("%s %s %d %s\n",
					v.Method, v.URI, v.Status, v.Latency)
			}
			return nil
		},
	}))
	e.Use(middleware.Recover())

	server.SetupRoutes(e, postHandler, userHandler, jwtSecret)

	e.Logger.Fatal(e.Start(":8080"))
}
