package models

import "gorm.io/gorm"

type Task struct {
    gorm.Model
    Title       string `json:"title"`
    Description string `json:"description"`
    AssignedTo  uint   `json:"assigned_to"`
    Status      string `json:"status"` // e.g., "pending", "completed"
}