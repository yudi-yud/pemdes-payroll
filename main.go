package main

import (
	"log"
	"pemdes-payroll/backend/config"
	"pemdes-payroll/backend/handlers"
	"pemdes-payroll/backend/middleware"
	"pemdes-payroll/backend/models"
	"pemdes-payroll/backend/repositories"
	"pemdes-payroll/backend/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Initialize database
	if err := config.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer config.CloseDB()

	// Auto migrate tables
	db := config.GetDB()
	err := db.AutoMigrate(
		&models.Jabatan{},
		&models.Karyawan{},
		&models.Gaji{},
		&models.Absensi{},
		&models.Lembur{},
		&models.User{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize repositories
	jabatanRepo := repositories.NewJabatanRepository(db)
	karyawanRepo := repositories.NewKaryawanRepository(db)
	gajiRepo := repositories.NewGajiRepository(db)
	laporanRepo := repositories.NewLaporanRepository(db)
	absensiRepo := repositories.NewAbsensiRepository(db)
	lemburRepo := repositories.NewLemburRepository(db)
	userRepo := repositories.NewUserRepository(db)

	// Initialize handlers
	jabatanHandler := handlers.NewJabatanHandler(jabatanRepo)
	karyawanHandler := handlers.NewKaryawanHandler(karyawanRepo)
	gajiHandler := handlers.NewGajiHandler(gajiRepo, karyawanRepo, lemburRepo)
	laporanHandler := handlers.NewLaporanHandler(laporanRepo, karyawanRepo, gajiRepo)
	absensiHandler := handlers.NewAbsensiHandler(absensiRepo, karyawanRepo)
	lemburHandler := handlers.NewLemburHandler(lemburRepo, karyawanRepo)
	authHandler := handlers.NewAuthHandler(userRepo)

	// Initialize default admin user
	if err := authHandler.InitAdmin(); err != nil {
		log.Printf("Warning: Failed to create default admin user: %v", err)
	} else {
		log.Println("Default admin user created: username=admin, password=admin123")
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Pemdes Payroll API",
		ServerHeader: "Pemdes Payroll",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(middleware.CORSConfig())

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Pemdes Payroll API is running",
		})
	})

	// Setup routes
	routes.SetupRoutes(app, jabatanHandler, karyawanHandler, gajiHandler, laporanHandler, absensiHandler, lemburHandler, authHandler)

	// Start server
	port := ":3000"
	log.Printf("Server starting on port %s", port)
	if err := app.Listen(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
