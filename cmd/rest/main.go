package main

import (
	"log"

	"github.com/daiyanuthsa/grpc-ecom-be/internal/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	// err := godotenv.Load("../../.env")
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
	godotenv.Load()

	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024,
	})
	app.Use(cors.New(cors.Config{
        // Izinkan origin frontend Anda
        AllowOrigins: "http://localhost:5173", 
        
        // Izinkan method yang diperlukan (GET, POST, OPTIONS, dll.)
        AllowMethods: "GET,POST,OPTIONS", 
        
        // Penting: Izinkan Authorization header karena Anda menggunakan JWT/Bearer token
        AllowHeaders: "Origin, Content-Type, Accept, Authorization", 
    }))

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "pong",
		})
	})

	app.Post("/product/upload", handler.UploadProductImageHandler)


	log.Println("Starting REST server on port 9000")
	if err := app.Listen(":9000"); err != nil {
		log.Fatalf("failed to run REST server: %v", err)
	}
}
