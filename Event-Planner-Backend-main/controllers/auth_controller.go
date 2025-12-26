package controllers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"event_planner_backend/config"
	"event_planner_backend/models"
	"event_planner_backend/utils"
	"gorm.io/gorm"
)

// Simple in-memory fallback store when DB is not connected.
var inMemoryUsers = map[string]*models.User{}

// SignupRequest represents the expected payload for signup.
type SignupRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=72"`
}

// LoginRequest represents the expected payload for login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Signup creates a new user with hashed password.
func Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid payload")
		return
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	name := strings.TrimSpace(req.Name)
	if name == "" || email == "" {
		utils.JSONError(c, http.StatusBadRequest, "name and email are required")
		return
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "failed to hash password")
		return
	}

    user := &models.User{Name: name, Email: email, PasswordHash: hash}

	if config.DB != nil {
		// Persist with GORM
		if err := config.DB.Create(user).Error; err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "duplicate") || strings.Contains(strings.ToLower(err.Error()), "unique") {
				utils.JSONError(c, http.StatusConflict, "email already registered")
				return
			}
			utils.JSONError(c, http.StatusInternalServerError, "failed to create user")
			return
		}
	} else {
		// In-memory fallback
		if _, exists := inMemoryUsers[email]; exists {
			utils.JSONError(c, http.StatusConflict, "email already registered")
			return
		}
		user.ID = uint(len(inMemoryUsers) + 1)
		inMemoryUsers[email] = user
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}

// findUserByEmail tries DB first, then in-memory.
func findUserByEmail(email string) (*models.User, error) {
	if config.DB != nil {
		var user models.User
		res := config.DB.Where("email = ?", email).First(&user)
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		if res.Error != nil {
			return nil, res.Error
		}
		return &user, nil
	}
	if u, ok := inMemoryUsers[email]; ok {
		return u, nil
	}
	return nil, nil
}

// Login verifies credentials and returns a JWT.
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid payload")
		return
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	user, err := findUserByEmail(email)
	if err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "failed to query user")
		return
	}
	if user == nil || !utils.CheckPassword(user.PasswordHash, req.Password) {
		utils.JSONError(c, http.StatusUnauthorized, "invalid email or password")
		return
	}

	secret := config.MustGetEnv("JWT_SECRET", "dev_secret_change_me")
	token, err := utils.GenerateJWT(secret, user.ID, user.Email)
	if err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "failed to generate token")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}
