package main

//go:generate buf generate

import (
	"log"
	"os"
	"study4cash/routes"
	"study4cash/websockets"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	if _, exists := os.LookupEnv("PRODUCTION"); !exists {
		godotenv.Load()
	}
	db := configDB()

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	// ROUTES
	routes.RouteUsers("/user", e, db)
	routes.RouteChats("/chats", e, db)
	e.GET("/ws", func(c *echo.Context) error { return websockets.WebsocketsHandler(c, db) })

	log.Println("Server is running on http://localhost:8080")
	if err := e.Start(":8080"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
