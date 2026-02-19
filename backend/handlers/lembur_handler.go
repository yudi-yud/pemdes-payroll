package handlers

import (
	"net/http"
	"pemdes-payroll/backend/models"
	"pemdes-payroll/backend/repositories"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type LemburHandler struct {
	lemburRepo    repositories.LemburRepository
	karyawanRepo repositories.KaryawanRepository
}

// NewLemburHandler creates a new Lembur handler
func NewLemburHandler(lemburRepo repositories.LemburRepository, karyawanRepo repositories.KaryawanRepository) *LemburHandler {
	return &LemburHandler{lemburRepo: lemburRepo, karyawanRepo: karyawanRepo}
}

// CreateLembur handles POST /api/lembur
func (h *LemburHandler) CreateLembur(c *fiber.Ctx) error {
	var req struct {
		KaryawanID uint   `json:"karyawan_id"`
		Tanggal    string `json:"tanggal"`
		JamMulai   string `json:"jam_mulai"`
		JamSelesai string `json:"jam_selesai"`
		TotalJam   float64 `json:"total_jam"`
		Keterangan string `json:"keterangan"`
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

	if req.TotalJam <= 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Total jam must be greater than 0",
		})
	}

	// Get karyawan data with jabatan to get overtime rate
	karyawan, err := h.karyawanRepo.GetByIDWithJabatan(req.KaryawanID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Karyawan not found",
		})
	}

	// Get tarif lembur from jabatan, default to 0 if no jabatan set
	tarifPerJam := 0.0
	if karyawan.Jabatan != nil {
		tarifPerJam = karyawan.Jabatan.TarifLemburPerJam
	}

	lembur := models.Lembur{
		KaryawanID:  req.KaryawanID,
		Tanggal:     tanggal,
		JamMulai:    req.JamMulai,
		JamSelesai:  req.JamSelesai,
		TotalJam:    req.TotalJam,
		TarifPerJam: tarifPerJam,
		TotalNominal: tarifPerJam * req.TotalJam,
		Keterangan:  req.Keterangan,
		Status:      "pending",
	}

	if err := h.lemburRepo.Create(&lembur); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create lembur",
		})
	}

	// Fetch with Karyawan data
	result, _ := h.lemburRepo.GetByID(lembur.ID)
	return c.Status(http.StatusCreated).JSON(result)
}

// GetAllLembur handles GET /api/lembur
func (h *LemburHandler) GetAllLembur(c *fiber.Ctx) error {
	lembur, err := h.lemburRepo.GetAll()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch lembur",
		})
	}

	return c.JSON(lembur)
}

// GetLemburByID handles GET /api/lembur/:id
func (h *LemburHandler) GetLemburByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	lembur, err := h.lemburRepo.GetByID(uint(id))
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Lembur not found",
		})
	}

	return c.JSON(lembur)
}

// GetLemburByPeriod handles GET /api/lembur/period?bulan=&tahun=
func (h *LemburHandler) GetLemburByPeriod(c *fiber.Ctx) error {
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

	lembur, err := h.lemburRepo.GetByPeriod(bulan, tahun)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch lembur",
		})
	}

	return c.JSON(lembur)
}

// GetLemburByKaryawan handles GET /api/lembur/karyawan/:id
func (h *LemburHandler) GetLemburByKaryawan(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	lembur, err := h.lemburRepo.GetByKaryawanID(uint(id))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch lembur",
		})
	}

	return c.JSON(lembur)
}

// UpdateLembur handles PUT /api/lembur/:id
func (h *LemburHandler) UpdateLembur(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	var req struct {
		Tanggal     string  `json:"tanggal"`
		JamMulai    string  `json:"jam_mulai"`
		JamSelesai  string  `json:"jam_selesai"`
		TotalJam    float64 `json:"total_jam"`
		Keterangan  string  `json:"keterangan"`
		Status      string  `json:"status"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Parse tanggal if provided
	var tanggal time.Time
	if req.Tanggal != "" {
		tanggal, err = time.Parse("2006-01-02", req.Tanggal)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid tanggal format. Use YYYY-MM-DD",
			})
		}
	}

	lembur := models.Lembur{
		Tanggal:    tanggal,
		JamMulai:   req.JamMulai,
		JamSelesai: req.JamSelesai,
		TotalJam:   req.TotalJam,
		Keterangan: req.Keterangan,
		Status:     req.Status,
	}

	if err := h.lemburRepo.Update(uint(id), &lembur); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update lembur",
		})
	}

	// Get updated data
	updated, _ := h.lemburRepo.GetByID(uint(id))
	return c.JSON(updated)
}

// DeleteLembur handles DELETE /api/lembur/:id
func (h *LemburHandler) DeleteLembur(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	if err := h.lemburRepo.Delete(uint(id)); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete lembur",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Lembur deleted successfully",
	})
}

// ApproveLembur handles PATCH /api/lembur/:id/approve
func (h *LemburHandler) ApproveLembur(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	var req struct {
		Status     string `json:"status"`
		ApproverID *uint  `json:"approver_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Status != "disetujui" && req.Status != "ditolak" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid status. Use 'disetujui' or 'ditolak'",
		})
	}

	// If approver_id is 0 or not provided, set to nil
	var approverID *uint
	if req.ApproverID != nil && *req.ApproverID > 0 {
		approverID = req.ApproverID
	}

	if err := h.lemburRepo.Approve(uint(id), approverID, req.Status); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to approve lembur",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Lembur status updated successfully",
	})
}

// RecalculateTarifLembur handles POST /api/lembur/recalculate-tarif
// This endpoint updates all lembur records with tarif_per_jam = 0 to use the correct tarif from jabatan
func (h *LemburHandler) RecalculateTarifLembur(c *fiber.Ctx) error {
	// Get all lembur with tarif = 0
	lemburList, err := h.lemburRepo.GetAll()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch lembur",
		})
	}

	updated := 0
	for _, lembur := range lemburList {
		if lembur.TarifPerJam == 0 {
			// Get karyawan with jabatan
			karyawan, err := h.karyawanRepo.GetByIDWithJabatan(lembur.KaryawanID)
			if err != nil {
				continue // Skip if karyawan not found
			}

			// Get tarif from jabatan
			tarifPerJam := 0.0
			if karyawan.Jabatan != nil {
				tarifPerJam = karyawan.Jabatan.TarifLemburPerJam
			}

			// Update lembur with new tarif and recalculate total
			lembur.TarifPerJam = tarifPerJam
			lembur.TotalNominal = tarifPerJam * lembur.TotalJam

			if err := h.lemburRepo.Update(lembur.ID, &lembur); err == nil {
				updated++
			}
		}
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Recalculation completed",
		"updated": updated,
		"total":   len(lemburList),
	})
}
