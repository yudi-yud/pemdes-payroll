package routes

import (
	"pemdes-payroll/backend/handlers"
	"pemdes-payroll/backend/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all API routes
func SetupRoutes(app *fiber.App,
	jabatanHandler *handlers.JabatanHandler,
	karyawanHandler *handlers.KaryawanHandler,
	gajiHandler *handlers.GajiHandler,
	laporanHandler *handlers.LaporanHandler,
	absensiHandler *handlers.AbsensiHandler,
	lemburHandler *handlers.LemburHandler,
	authHandler *handlers.AuthHandler,
) {
	// Public routes (no auth required)
	app.Post("/api/auth/login", authHandler.Login)

	// API routes with auth middleware
	api := app.Group("/api", middleware.AuthMiddleware)

	// Auth routes
	api.Get("/auth/me", authHandler.Me)
	api.Get("/auth/users", authHandler.GetUsers)
	api.Post("/auth/users", authHandler.CreateUser)
	api.Put("/auth/users/:id", authHandler.UpdateUser)
	api.Delete("/auth/users/:id", authHandler.DeleteUser)
	api.Patch("/auth/users/:id/toggle", authHandler.ToggleUserActive)
	api.Put("/auth/change-password", authHandler.ChangePassword)

	// Jabatan routes
	api.Get("/jabatan", jabatanHandler.GetAllJabatan)
	api.Get("/jabatan/:id", jabatanHandler.GetJabatanByID)
	api.Post("/jabatan", jabatanHandler.CreateJabatan)
	api.Put("/jabatan/:id", jabatanHandler.UpdateJabatan)
	api.Delete("/jabatan/:id", jabatanHandler.DeleteJabatan)

	// Karyawan routes
	api.Get("/karyawan", karyawanHandler.GetAllKaryawan)
	api.Get("/karyawan/search", karyawanHandler.SearchKaryawan)
	api.Get("/karyawan/:id", karyawanHandler.GetKaryawanByID)
	api.Post("/karyawan", karyawanHandler.CreateKaryawan)
	api.Put("/karyawan/:id", karyawanHandler.UpdateKaryawan)
	api.Delete("/karyawan/:id", karyawanHandler.DeleteKaryawan)
	api.Get("/karyawan/status/:status", karyawanHandler.GetActiveKaryawan)

	// Absensi routes
	api.Get("/absensi", absensiHandler.GetAllAbsensi)
	api.Get("/absensi/:id", absensiHandler.GetAbsensiByID)
	api.Get("/absensi/karyawan/:id", absensiHandler.GetAbsensiByKaryawan)
	api.Get("/absensi/rekap/:karyawan_id", absensiHandler.GetRekapAbsensi)
	api.Post("/absensi", absensiHandler.CreateAbsensi)
	api.Put("/absensi/:id", absensiHandler.UpdateAbsensi)
	api.Delete("/absensi/:id", absensiHandler.DeleteAbsensi)
	api.Get("/absensi/export/excel", absensiHandler.ExportAbsensiExcel)
	api.Get("/absensi/export/karyawan/:id/pdf", absensiHandler.ExportAbsensiPDF)

	// Lembur routes
	api.Get("/lembur", lemburHandler.GetAllLembur)
	api.Get("/lembur/period", lemburHandler.GetLemburByPeriod)
	api.Get("/lembur/:id", lemburHandler.GetLemburByID)
	api.Get("/lembur/karyawan/:id", lemburHandler.GetLemburByKaryawan)
	api.Post("/lembur", lemburHandler.CreateLembur)
	api.Put("/lembur/:id", lemburHandler.UpdateLembur)
	api.Delete("/lembur/:id", lemburHandler.DeleteLembur)
	api.Patch("/lembur/:id/approve", lemburHandler.ApproveLembur)
	api.Post("/lembur/recalculate-tarif", lemburHandler.RecalculateTarifLembur)

	// Gaji routes - static routes first, then parameterized routes
	api.Get("/gaji", gajiHandler.GetAllGaji)
	api.Get("/gaji/period", gajiHandler.GetGajiByPeriod)
	api.Get("/gaji/my-slips", gajiHandler.GetMySlipGaji)  // Static route before :id
	api.Get("/gaji/generate-batch", gajiHandler.GenerateBatch)  // Static route
	api.Get("/gaji/:id", gajiHandler.GetGajiByID)
	api.Get("/gaji/karyawan/:id", gajiHandler.GetGajiByKaryawanID)
	api.Get("/gaji/slip/:id", gajiHandler.GetSlipGaji)
	api.Post("/gaji", gajiHandler.CreateGaji)
	api.Post("/gaji/generate-batch", gajiHandler.GenerateBatch)
	api.Put("/gaji/:id", gajiHandler.UpdateGaji)
	api.Delete("/gaji/:id", gajiHandler.DeleteGaji)
	api.Patch("/gaji/:id/status", gajiHandler.UpdateGajiStatus)

	// Laporan routes
	api.Get("/laporan/gaji", laporanHandler.GetLaporanGajiByPeriod)
	api.Get("/laporan/gaji/karyawan/:id", laporanHandler.GetRiwayatGajiKaryawan)
	api.Get("/laporan/rekap", laporanHandler.GetRekapGaji)

	// Export routes
	api.Get("/laporan/export/excel", laporanHandler.ExportLaporanExcel)
	api.Get("/laporan/export/karyawan/:id/pdf", laporanHandler.ExportKaryawanPDF)
}
