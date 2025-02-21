package controllers

import (
    "net/http"
    "ai-task-manager-backend/models"
    "ai-task-manager-backend/utils"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func Register(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var user models.User
        if err := c.ShouldBindJSON(&user); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        db.Create(&user)
        c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
    }
}

func Login(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var input struct {
            Username string `json:"username"`
            Password string `json:"password"`
        }
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        var user models.User
        if err := db.Where("username = ?", input.Username).First(&user).Error; err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
            return
        }
        token, _ := utils.GenerateToken(uint(user.ID))
        c.JSON(http.StatusOK, gin.H{"token": token})
    }
}