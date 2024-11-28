package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

type Cred struct {
	Username string
	Password string
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// var cred Cred
	// cred.Username = os.Getenv("APP_USERNAME")
	// cred.Password = os.Getenv("APP_PASSWORD")

	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())
	app.Use(cors.New(cors.Config{
		AllowMethods: "GET,POST",
	}))

	app.Use(func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		log.Println((authHeader[7:]), os.Getenv("APP_PASSWORD"), len(authHeader[7:]), len(os.Getenv("APP_PASSWORD")), authHeader[7:] == os.Getenv("APP_PASSWORD"), authHeader[:7] == "Bearer ")
		if authHeader[7:] == os.Getenv("APP_PASSWORD") && authHeader[:7] == "Bearer " {
			token := authHeader[7:]
			fmt.Println("Bearer Token:", token)
			return c.Next()
		}
		return c.Status(500).JSON(&fiber.Map{
			"mess": "something went wrong",
		})
	})

	app.Post("/query-execute", QueryExecute)
	app.Get("/health-check", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(&fiber.Map{
			"response": "All Ok",
		})
	})

	log.Fatal(app.Listen(":7777"))
	defer CloseDB()

}
