package controllers

import (
    "net/http"
    "ai-task-manager-backend/models"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func CreateTask(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var task models.Task
        if err := c.ShouldBindJSON(&task); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        db.Create(&task)
        c.JSON(http.StatusCreated, task)
    }
}

func GetTasks(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var tasks []models.Task
        db.Find(&tasks)
        c.JSON(http.StatusOK, tasks)
    }
}