package handlers

import (
	"net/http"
	"pemdes-payroll/backend/middleware"
	"pemdes-payroll/backend/models"
	"pemdes-payroll/backend/repositories"
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
	expirationTime := time.Now().Add(24 * time.Hour)
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
