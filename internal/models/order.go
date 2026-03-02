package models

import (
	"time"

	"gorm.io/gorm"
)

// Order
type Order struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	CustomerName string         `json:"customer_name" binding:"required"`
	Amount       float64        `json:"amount" binding:"required,gt=0"`
	Status       string         `gorm:"default:PENDING" json:"status"` // PENDING, PROCESSING, COMPLETED
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete  (not shown in JSON)
}
