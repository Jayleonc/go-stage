package main

import (
	"github.com/Jayleonc/go-stage/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r := gin.Default()

	r.Use(gin.Logger())

	routes.AuthRoutes(r)
	routes.UserRoutes(r)

	err = r.Run(":" + port)
	if err != nil {
		return
	}
}
