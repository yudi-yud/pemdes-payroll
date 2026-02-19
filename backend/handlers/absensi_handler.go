package handlers

import (
	"fmt"
	"net/http"
	"pemdes-payroll/backend/models"
	"pemdes-payroll/backend/repositories"
	"pemdes-payroll/backend/services"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AbsensiHandler struct {
	absensiRepo   repositories.AbsensiRepository
	karyawanRepo  repositories.KaryawanRepository
	exportService *services.ExportService
}

// NewAbsensiHandler creates a new Absensi handler
func NewAbsensiHandler(absensiRepo repositories.AbsensiRepository, karyawanRepo repositories.KaryawanRepository) *AbsensiHandler {
	return &AbsensiHandler{
		absensiRepo:   absensiRepo,
		karyawanRepo:  karyawanRepo,
		exportService: services.NewExportService(),
	}
}

// CreateAbsensi handles POST /api/absensi
func (h *AbsensiHandler) CreateAbsensi(c *fiber.Ctx) error {
	var req struct {
		KaryawanID uint           `json:"karyawan_id"`
		Tanggal    string         `json:"tanggal"`
		JamMasuk   string         `json:"jam_masuk"`
		JamKeluar  string         `json:"jam_keluar"`
		Status     models.AbsensiStatus `json:"status"`
		Keterangan string         `json:"keterangan"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validation
	if req.KaryawanID == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Karyawan ID is required",
		})
	}
	if req.Tanggal == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Tanggal is required",
		})
	}

	// Parse tanggal
	tanggal, err := time.Parse("2006-01-02", req.Tanggal)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tanggal format. Use YYYY-MM-DD",
		})
	}

	absensi := models.Absensi{
		KaryawanID: req.KaryawanID,
		Tanggal:    tanggal,
		JamMasuk:   req.JamMasuk,
		JamKeluar:  req.JamKeluar,
		Status:     req.Status,
		Keterangan: req.Keterangan,
	}

	if absensi.Status == "" {
		absensi.Status = models.AbsensiHadir
	}

	if err := h.absensiRepo.Create(&absensi); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create absensi",
		})
	}

	// Fetch with Karyawan data
	result, _ := h.absensiRepo.GetByID(absensi.ID)
	return c.Status(http.StatusCreated).JSON(result)
}

// GetAllAbsensi handles GET /api/absensi
func (h *AbsensiHandler) GetAllAbsensi(c *fiber.Ctx) error {
	absensi, err := h.absensiRepo.GetAll()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch absensi",
		})
	}

	return c.JSON(absensi)
}

// GetAbsensiByID handles GET /api/absensi/:id
func (h *AbsensiHandler) GetAbsensiByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	absensi, err := h.absensiRepo.GetByID(uint(id))
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Absensi not found",
		})
	}

	return c.JSON(absensi)
}

// GetAbsensiByKaryawan handles GET /api/absensi/karyawan/:id
func (h *AbsensiHandler) GetAbsensiByKaryawan(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time
	var errDate error

	if startDateStr != "" && endDateStr != "" {
		startDate, errDate = time.Parse("2006-01-02", startDateStr)
		if errDate != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid start_date format. Use YYYY-MM-DD",
			})
		}
		endDate, errDate = time.Parse("2006-01-02", endDateStr)
		if errDate != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid end_date format. Use YYYY-MM-DD",
			})
		}
	} else {
		// Default to current month
		now := time.Now()
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		endDate = startDate.AddDate(0, 1, -1)
	}

	absensi, err := h.absensiRepo.GetByKaryawanID(uint(id), startDate, endDate)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch absensi",
		})
	}

	return c.JSON(absensi)
}

// GetRekapAbsensi handles GET /api/absensi/rekap/:karyawan_id?bulan=&tahun=
func (h *AbsensiHandler) GetRekapAbsensi(c *fiber.Ctx) error {
	karyawanID, err := strconv.ParseUint(c.Params("karyawan_id"), 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Karyawan ID",
		})
	}

	bulan, _ := strconv.Atoi(c.Query("bulan", strconv.Itoa(int(time.Now().Month()))))
	tahun, _ := strconv.Atoi(c.Query("tahun", strconv.Itoa(time.Now().Year())))

	if bulan < 1 || bulan > 12 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid bulan parameter",
		})
	}

	rekap, err := h.absensiRepo.GetRekapBulanan(uint(karyawanID), bulan, tahun)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get rekap absensi",
		})
	}

	return c.JSON(rekap)
}

// UpdateAbsensi handles PUT /api/absensi/:id
func (h *AbsensiHandler) UpdateAbsensi(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	var req struct {
		JamMasuk   string         `json:"jam_masuk"`
		JamKeluar  string         `json:"jam_keluar"`
		Status     models.AbsensiStatus `json:"status"`
		Keterangan string         `json:"keterangan"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	absensi := models.Absensi{
		JamMasuk:   req.JamMasuk,
		JamKeluar:  req.JamKeluar,
		Status:     req.Status,
		Keterangan: req.Keterangan,
	}

	if err := h.absensiRepo.Update(uint(id), &absensi); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update absensi",
		})
	}

	// Get updated data
	updated, _ := h.absensiRepo.GetByID(uint(id))
	return c.JSON(updated)
}

// DeleteAbsensi handles DELETE /api/absensi/:id
func (h *AbsensiHandler) DeleteAbsensi(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	if err := h.absensiRepo.Delete(uint(id)); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete absensi",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Absensi deleted successfully",
	})
}

// ExportAbsensiExcel handles GET /api/absensi/export/excel
func (h *AbsensiHandler) ExportAbsensiExcel(c *fiber.Ctx) error {
	// Get all absensi
	absensiList, err := h.absensiRepo.GetAll()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch absensi",
		})
	}

	// Generate Excel
	data, err := h.exportService.ExportAbsensiToExcel(absensiList)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate Excel",
		})
	}

	// Set headers for download
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename=absensi_karyawan.xlsx")

	return c.Send(data)
}

// ExportAbsensiPDF handles GET /api/absensi/export/karyawan/:id/pdf?bulan=&tahun=
func (h *AbsensiHandler) ExportAbsensiPDF(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Karyawan ID",
		})
	}

	bulan, _ := strconv.Atoi(c.Query("bulan", strconv.Itoa(int(time.Now().Month()))))
	tahun, _ := strconv.Atoi(c.Query("tahun", strconv.Itoa(time.Now().Year())))

	if bulan < 1 || bulan > 12 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid bulan parameter",
		})
	}

	// Get karyawan
	karyawan, err := h.karyawanRepo.GetByIDWithJabatan(uint(id))
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Karyawan not found",
		})
	}

	// Get absensi for the period
	startDate := time.Date(tahun, time.Month(bulan), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1)
	absensiList, err := h.absensiRepo.GetByKaryawanID(uint(id), startDate, endDate)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch absensi",
		})
	}

	// Get rekap
	rekap, err := h.absensiRepo.GetRekapBulanan(uint(id), bulan, tahun)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get rekap",
		})
	}

	// Generate PDF
	data, err := h.exportService.ExportAbsensiToPDF(karyawan, absensiList, rekap, bulan, tahun)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate PDF",
		})
	}

	// Set headers for download
	c.Set("Content-Type", "application/pdf")
	filename := fmt.Sprintf("absensi_%s_%s%d.pdf", karyawan.Nama, getMonthName(bulan), tahun)
	c.Set("Content-Disposition", "attachment; filename="+filename)

	return c.Send(data)
}
