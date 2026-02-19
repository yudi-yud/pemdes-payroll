package handlers

import (
	"fmt"
	"net/http"
	"pemdes-payroll/backend/repositories"
	"pemdes-payroll/backend/services"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type LaporanHandler struct {
	laporanRepo  repositories.LaporanRepository
	karyawanRepo repositories.KaryawanRepository
	gajiRepo     repositories.GajiRepository
	exportSvc    *services.ExportService
}

// NewLaporanHandler creates a new Laporan handler
func NewLaporanHandler(
	laporanRepo repositories.LaporanRepository,
	karyawanRepo repositories.KaryawanRepository,
	gajiRepo repositories.GajiRepository,
) *LaporanHandler {
	return &LaporanHandler{
		laporanRepo:  laporanRepo,
		karyawanRepo: karyawanRepo,
		gajiRepo:     gajiRepo,
		exportSvc:    services.NewExportService(),
	}
}

// GetLaporanGajiByPeriod handles GET /api/laporan/gaji?bulan=&tahun=
func (h *LaporanHandler) GetLaporanGajiByPeriod(c *fiber.Ctx) error {
	bulan, _ := strconv.Atoi(c.Query("bulan", "0"))
	tahun, _ := strconv.Atoi(c.Query("tahun", "0"))

	if bulan < 1 || bulan > 12 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid bulan parameter",
		})
	}
	if tahun < 2000 || tahun > 2100 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tahun parameter",
		})
	}

	laporan, err := h.laporanRepo.GetLaporanGajiByPeriod(bulan, tahun)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch laporan",
		})
	}

	return c.JSON(laporan)
}

// GetRiwayatGajiKaryawan handles GET /api/laporan/gaji/karyawan/:id
func (h *LaporanHandler) GetRiwayatGajiKaryawan(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	laporan, err := h.laporanRepo.GetRiwayatGajiKaryawan(uint(id))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch riwayat gaji",
		})
	}

	return c.JSON(laporan)
}

// GetRekapGaji handles GET /api/laporan/rekap?bulan=&tahun=
func (h *LaporanHandler) GetRekapGaji(c *fiber.Ctx) error {
	bulan, _ := strconv.Atoi(c.Query("bulan", "0"))
	tahun, _ := strconv.Atoi(c.Query("tahun", "0"))

	if bulan < 1 || bulan > 12 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid bulan parameter",
		})
	}
	if tahun < 2000 || tahun > 2100 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tahun parameter",
		})
	}

	rekap, err := h.laporanRepo.GetRekapGaji(bulan, tahun)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch rekap gaji",
		})
	}

	return c.JSON(rekap)
}

// ExportLaporanExcel handles GET /api/laporan/export/excel?bulan=&tahun=
func (h *LaporanHandler) ExportLaporanExcel(c *fiber.Ctx) error {
	bulan, _ := strconv.Atoi(c.Query("bulan", "0"))
	tahun, _ := strconv.Atoi(c.Query("tahun", "0"))

	if bulan < 1 || bulan > 12 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid bulan parameter",
		})
	}
	if tahun < 2000 || tahun > 2100 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tahun parameter",
		})
	}

	// Get laporan data for the period
	laporanList, err := h.laporanRepo.GetLaporanGajiByPeriod(bulan, tahun)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch laporan data",
		})
	}

	// Generate Excel
	data, err := h.exportSvc.ExportToExcel(laporanList, bulan, tahun)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate Excel",
		})
	}

	// Set headers for download
	monthName := getMonthName(bulan)
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=Laporan_Gaji_%s_%d.xlsx", monthName, tahun))

	return c.Send(data)
}

// getMonthName returns Indonesian month name
func getMonthName(month int) string {
	months := []string{
		"Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}
	if month >= 1 && month <= 12 {
		return months[month-1]
	}
	return ""
}

// ExportKaryawanPDF handles GET /api/laporan/export/karyawan/:id/pdf
func (h *LaporanHandler) ExportKaryawanPDF(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	// Get karyawan with jabatan
	karyawan, err := h.karyawanRepo.GetByIDWithJabatan(uint(id))
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Karyawan not found",
		})
	}

	// Get all gaji for this karyawan
	gajiList, err := h.gajiRepo.GetByKaryawanID(uint(id))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch gaji data",
		})
	}

	// Generate PDF
	data, err := h.exportSvc.ExportToPDF(karyawan, gajiList)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate PDF",
		})
	}

	// Set headers for download
	filename := strings.ReplaceAll(karyawan.Nama, " ", "_")
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=Laporan_Gaji_%s.pdf", filename))

	return c.Send(data)
}
