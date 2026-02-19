package handlers

import (
	"net/http"
	"pemdes-payroll/backend/models"
	"pemdes-payroll/backend/repositories"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type JabatanHandler struct {
	repo repositories.JabatanRepository
}

// NewJabatanHandler creates a new Jabatan handler
func NewJabatanHandler(repo repositories.JabatanRepository) *JabatanHandler {
	return &JabatanHandler{repo: repo}
}

// CreateJabatan handles POST /api/jabatan
func (h *JabatanHandler) CreateJabatan(c *fiber.Ctx) error {
	var req models.Jabatan
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validation
	if req.NamaJabatan == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Nama jabatan is required",
		})
	}
	if req.GajiPokok <= 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Gaji pokok must be greater than 0",
		})
	}

	if err := h.repo.Create(&req); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create jabatan",
		})
	}

	return c.Status(http.StatusCreated).JSON(req)
}

// GetAllJabatan handles GET /api/jabatan
func (h *JabatanHandler) GetAllJabatan(c *fiber.Ctx) error {
	jabatans, err := h.repo.GetAll()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch jabatans",
		})
	}

	return c.JSON(jabatans)
}

// GetJabatanByID handles GET /api/jabatan/:id
func (h *JabatanHandler) GetJabatanByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	jabatan, err := h.repo.GetByID(uint(id))
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Jabatan not found",
		})
	}

	return c.JSON(jabatan)
}

// UpdateJabatan handles PUT /api/jabatan/:id
func (h *JabatanHandler) UpdateJabatan(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	var req models.Jabatan
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validation
	if req.NamaJabatan == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Nama jabatan is required",
		})
	}
	if req.GajiPokok <= 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Gaji pokok must be greater than 0",
		})
	}

	if err := h.repo.Update(uint(id), &req); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update jabatan",
		})
	}

	// Get updated data
	updated, _ := h.repo.GetByID(uint(id))
	return c.JSON(updated)
}

// DeleteJabatan handles DELETE /api/jabatan/:id
func (h *JabatanHandler) DeleteJabatan(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete jabatan",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Jabatan deleted successfully",
	})
}
