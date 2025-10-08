package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
	"workout-tracker/models"
	"workout-tracker/utils"
)

type AuthController struct{ DB *gorm.DB }

type signupReq struct {
	Name     string `json:"name" binding:"required,min=2"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// @Summary Sign up
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body signupReq true "Signup"
// @Success 201
// @Router /auth/signup [post]
func (a *AuthController) Signup(c *gin.Context) {
	var req signupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hash, _ := utils.HashPassword(req.Password)
	u := models.User{Name: req.Name, Email: req.Email, PasswordHash: hash}
	if err := a.DB.Create(&u).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email may already exist"})
		return
	}
	c.Status(http.StatusCreated)
}

// @Summary Login
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body loginReq true "Login"
// @Success 200 {object} map[string]string
// @Router /auth/login [post]
func (a *AuthController) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user models.User
	if err := a.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if !utils.CheckPassword(user.PasswordHash, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	tok, _ := utils.GenerateToken(user.ID, user.Email, 24*time.Hour)
	c.JSON(http.StatusOK, gin.H{"token": tok})
}
