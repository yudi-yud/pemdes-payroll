package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// CORSConfig returns CORS middleware configuration
func CORSConfig() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173,http://localhost:3000,http://76.13.23.109:8083,http://76.13.23.109:3001",
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: false,
		MaxAge:           3600,
	})
}
