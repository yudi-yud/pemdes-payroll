package handlers

import (
	"net/http"
	"pemdes-payroll/backend/models"
	"pemdes-payroll/backend/repositories"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type KaryawanHandler struct {
	repo repositories.KaryawanRepository
}

// NewKaryawanHandler creates a new Karyawan handler
func NewKaryawanHandler(repo repositories.KaryawanRepository) *KaryawanHandler {
	return &KaryawanHandler{repo: repo}
}

// CreateKaryawan handles POST /api/karyawan
func (h *KaryawanHandler) CreateKaryawan(c *fiber.Ctx) error {
	var req struct {
		NIK             string     `json:"nik"`
		Nama            string     `json:"nama"`
		Email           string     `json:"email"`
		Telepon         string     `json:"telepon"`
		Alamat          string     `json:"alamat"`
		JabatanID       *uint      `json:"jabatan_id"`
		TanggalBergabung *string   `json:"tanggal_bergabung"`
		Status          models.KaryawanStatus `json:"status"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validation
	if req.NIK == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "NIK is required",
		})
	}
	if req.Nama == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Nama is required",
		})
	}

	karyawan := models.Karyawan{
		NIK:       req.NIK,
		Nama:      req.Nama,
		Email:     req.Email,
		Telepon:   req.Telepon,
		Alamat:    req.Alamat,
		JabatanID: req.JabatanID,
		Status:    req.Status,
	}

	if karyawan.Status == "" {
		karyawan.Status = models.StatusAktif
	}

	if req.TanggalBergabung != nil && *req.TanggalBergabung != "" {
		parsedDate, err := time.Parse("2006-01-02", *req.TanggalBergabung)
		if err == nil {
			karyawan.TanggalBergabung = &parsedDate
		}
	}

	if err := h.repo.Create(&karyawan); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create karyawan",
		})
	}

	// Fetch with Jabatan data
	result, _ := h.repo.GetByIDWithJabatan(karyawan.ID)
	return c.Status(http.StatusCreated).JSON(result)
}

// GetAllKaryawan handles GET /api/karyawan
func (h *KaryawanHandler) GetAllKaryawan(c *fiber.Ctx) error {
	karyawan, err := h.repo.GetAll()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch karyawan",
		})
	}

	return c.JSON(karyawan)
}

// SearchKaryawan handles GET /api/karyawan/search?q=
func (h *KaryawanHandler) SearchKaryawan(c *fiber.Ctx) error {
	keyword := c.Query("q", "")
	karyawan, err := h.repo.Search(keyword)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search karyawan",
		})
	}

	return c.JSON(karyawan)
}

// GetKaryawanByID handles GET /api/karyawan/:id
func (h *KaryawanHandler) GetKaryawanByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	karyawan, err := h.repo.GetByIDWithJabatan(uint(id))
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Karyawan not found",
		})
	}

	return c.JSON(karyawan)
}

// UpdateKaryawan handles PUT /api/karyawan/:id
func (h *KaryawanHandler) UpdateKaryawan(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	var req struct {
		NIK             string     `json:"nik"`
		Nama            string     `json:"nama"`
		Email           string     `json:"email"`
		Telepon         string     `json:"telepon"`
		Alamat          string     `json:"alamat"`
		JabatanID       *uint      `json:"jabatan_id"`
		TanggalBergabung *string   `json:"tanggal_bergabung"`
		Status          models.KaryawanStatus `json:"status"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validation
	if req.NIK == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "NIK is required",
		})
	}
	if req.Nama == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Nama is required",
		})
	}

	karyawan := models.Karyawan{
		NIK:       req.NIK,
		Nama:      req.Nama,
		Email:     req.Email,
		Telepon:   req.Telepon,
		Alamat:    req.Alamat,
		JabatanID: req.JabatanID,
		Status:    req.Status,
	}

	if req.TanggalBergabung != nil && *req.TanggalBergabung != "" {
		parsedDate, err := time.Parse("2006-01-02", *req.TanggalBergabung)
		if err == nil {
			karyawan.TanggalBergabung = &parsedDate
		}
	}

	if err := h.repo.Update(uint(id), &karyawan); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update karyawan",
		})
	}

	// Get updated data
	updated, _ := h.repo.GetByIDWithJabatan(uint(id))
	return c.JSON(updated)
}

// DeleteKaryawan handles DELETE /api/karyawan/:id
func (h *KaryawanHandler) DeleteKaryawan(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete karyawan",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Karyawan deleted successfully",
	})
}

// GetActiveKaryawan handles GET /api/karyawan/status/:status
func (h *KaryawanHandler) GetActiveKaryawan(c *fiber.Ctx) error {
	status := models.KaryawanStatus(c.Params("status"))
	if status != models.StatusAktif && status != models.StatusNonAktif {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid status. Use 'aktif' or 'non_aktif'",
		})
	}

	karyawan, err := h.repo.GetByStatus(status)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch karyawan",
		})
	}

	return c.JSON(karyawan)
}
