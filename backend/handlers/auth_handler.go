package handlers

import (
	"net/http"
	"pemdes-payroll/backend/middleware"
	"pemdes-payroll/backend/models"
	"pemdes-payroll/backend/repositories"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("pemdes-payroll-secret-key-2024")

type AuthHandler struct {
	userRepo repositories.UserRepository
}

// NewAuthHandler creates a new Auth handler
func NewAuthHandler(userRepo repositories.UserRepository) *AuthHandler {
	return &AuthHandler{userRepo: userRepo}
}

// LoginRequest represents login request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents login response
type LoginResponse struct {
	Token string       `json:"token"`
	User  models.User `json:"user"`
}

// Login handles POST /api/auth/login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validation
	if req.Username == "" || req.Password == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Username and password are required",
		})
	}

	// Get user
	user, err := h.userRepo.GetByUsername(req.Username)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid username or password",
		})
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid username or password",
		})
	}

	// Create JWT token
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &middleware.Claims{
		UserID:   user.ID,
		Username: user.Username,
		Name:     user.Name,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	// Remove password from response
	user.Password = ""

	return c.JSON(LoginResponse{
		Token: tokenString,
		User:  *user,
	})
}

// Me handles GET /api/auth/me
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	user := c.Locals("user").(*middleware.Claims)

	userData, err := h.userRepo.GetByID(user.UserID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Remove password from response
	userData.Password = ""

	return c.JSON(userData)
}

// CreateUserRequest represents create user request
type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

// CreateUser handles POST /api/auth/users (create new user)
func (h *AuthHandler) CreateUser(c *fiber.Ctx) error {
	var req CreateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validation
	if req.Username == "" || req.Password == "" || req.Name == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Username, password, and name are required",
		})
	}

	// Check if username already exists
	_, err := h.userRepo.GetByUsername(req.Username)
	if err == nil {
		return c.Status(http.StatusConflict).JSON(fiber.Map{
			"error": "Username already exists",
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	user := models.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Name:     req.Name,
		Email:    req.Email,
		Role:     req.Role,
		IsActive: true,
	}

	if user.Role == "" {
		user.Role = "admin"
	}

	if err := h.userRepo.Create(&user); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	// Remove password from response
	user.Password = ""

	return c.Status(http.StatusCreated).JSON(user)
}

// InitAdmin creates default admin user if not exists
func (h *AuthHandler) InitAdmin() error {
	// Check if admin exists
	_, err := h.userRepo.GetByUsername("admin")
	if err == nil {
		return nil // Admin already exists
	}

	// Create default admin
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)

	user := models.User{
		Username: "admin",
		Password: string(hashedPassword),
		Name:     "Administrator",
		Email:    "admin@pemdes.desa",
		Role:     "admin",
		IsActive: true,
	}

	return h.userRepo.Create(&user)
}

// GetUsers handles GET /api/auth/users - Get all users (admin only)
func (h *AuthHandler) GetUsers(c *fiber.Ctx) error {
	users, err := h.userRepo.GetAll()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}

	// Remove passwords from response
	for i := range users {
		users[i].Password = ""
	}

	return c.JSON(users)
}

// UpdateUserRequest represents update user request
type UpdateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	IsActive *bool  `json:"is_active"`
}

// UpdateUser handles PUT /api/auth/users/:id - Update user
func (h *AuthHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get existing user
	existing, err := h.userRepo.GetByID(uint(id))
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Update fields
	existing.Name = req.Name
	existing.Email = req.Email
	if req.Role != "" {
		existing.Role = req.Role
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}

	if err := h.userRepo.Update(uint(id), existing); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	// Get updated user
	updated, _ := h.userRepo.GetByID(uint(id))
	updated.Password = ""

	return c.JSON(updated)
}

// DeleteUser handles DELETE /api/auth/users/:id - Delete user
func (h *AuthHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	// Prevent deleting yourself
	user := c.Locals("user").(*middleware.Claims)
	if uint(id) == user.UserID {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"error": "Cannot delete your own account",
		})
	}

	if err := h.userRepo.Delete(uint(id)); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}

// ToggleUserActive handles PATCH /api/auth/users/:id/toggle - Toggle user active status
func (h *AuthHandler) ToggleUserActive(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	// Prevent toggling yourself
	user := c.Locals("user").(*middleware.Claims)
	if uint(id) == user.UserID {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"error": "Cannot toggle your own account",
		})
	}

	if err := h.userRepo.ToggleActive(uint(id)); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to toggle user status",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "User status toggled successfully",
	})
}

// ChangePasswordRequest represents change password request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// ChangePassword handles PUT /api/auth/change-password - Change user password
func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	user := c.Locals("user").(*middleware.Claims)

	var req ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get user
	userData, err := h.userRepo.GetByID(user.UserID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(req.OldPassword)); err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Old password is incorrect",
		})
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	userData.Password = string(hashedPassword)
	if err := h.userRepo.Update(user.UserID, userData); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update password",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Password changed successfully",
	})
}
