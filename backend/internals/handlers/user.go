package handlers

import (
	"aetrix/observer/internals/lib"
	"aetrix/observer/internals/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserHandler handles all user-related HTTP requests and database operations.
// It contains a database connection for performing CRUD operations on user entities.
type UserHandler struct {
	DB *gorm.DB
}

// NewUserHandler creates and returns a new instance of UserHandler with the provided database connection.
// This is the constructor function for UserHandler.
func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{DB: db}
}

// SignupRequest represents the request payload for user registration.
// All fields are required and email must be in valid email format.
type SignupRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

// LoginRequest represents the request payload for user authentication.
// Both username and password are required fields.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UpdateUserRequest represents the request payload for updating user information.
// Both fields are optional, but if email is provided, it must be in valid format.
type UpdateUserRequest struct {
	Username string `json:"username" binding:"omitempty"`
	Email    string `json:"email" binding:"omitempty,email"`
}

// ChangePasswordRequest represents the request payload for changing user password.
// Both old and new passwords are required for security verification.
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// Signup handles user registration requests.
// It validates the request data, creates a new user in the database,
// and returns a success response with the created user ID.
// Returns 400 for validation errors or 201 for successful creation.
func (h *UserHandler) Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	value := h.DB.Create(&models.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	})

	c.JSON(201, gin.H{"message": "User created successfully", "user_id": value.RowsAffected})
}

// GetUser retrieves the authenticated user's information.
// It extracts the user ID from the context (set by authentication middleware),
// fetches the user from database, and returns the user details without password.
// Returns 401 for unauthenticated requests or 404 if user not found.
func (h *UserHandler) GetUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if err := h.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	c.JSON(200, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	})
}

// Login handles user authentication requests.
// It validates credentials against the database and generates a JWT token
// for successful authentication. Returns the user ID and token on success.
// Returns 400 for validation errors, 401 for invalid credentials, or 500 for token generation errors.
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.DB.First(&user, "username = ? AND password = ?", req.Username, req.Password).Error; err != nil {
		c.JSON(401, gin.H{"error": "Invalid username or password"})
		return
	}

	token, err := lib.GenerateToken(user.ID, lib.USER_JWT_SECRET) // Implement this function to generate JWT token
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}
	c.JSON(200, gin.H{"message": "Login successful", "user_id": user.ID, "token": token})
}

// UpdateUserDetailas handles requests to update user information.
// It allows authenticated users to update their username and/or email.
// Only provided fields will be updated, empty fields are ignored.
// Returns 401 for unauthenticated requests, 404 if user not found, or 200 for success.
func (h *UserHandler) UpdateUserDetailas(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	h.DB.Save(&user)
	c.JSON(200, gin.H{"message": "User updated successfully"})
}

// ChangePassword handles password change requests for authenticated users.
// It verifies the old password before updating to the new password for security.
// Requires both old and new passwords in the request.
// Returns 401 for unauthenticated requests, 400 for incorrect old password, 404 if user not found, or 200 for success.
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	if user.Password != req.OldPassword {
		c.JSON(400, gin.H{"error": "Old password is incorrect"})
		return
	}

	user.Password = req.NewPassword
	h.DB.Save(&user)
	c.JSON(200, gin.H{"message": "Password changed successfully"})
}
