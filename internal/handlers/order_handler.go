package handlers

import (
	"logistics-simulator/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OrderHandler store dependency cho handler (DB, JobChan)
type OrderHandler struct {
	DB      *gorm.DB
	JobChan chan<- uint // send-only channel to send order IDs to workers
}

// CreateOrder get data from request, create new order, save to DB and send order ID to JobChan for worker to process
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var input struct {
		CustomerName string  `json:"customer_name" binding:"required"`
		Amount       float64 `json:"amount" binding:"required,gt=0"`
	}

	// 1. form validation - ensure we have all required fields and valid data
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. create new order instance with status PENDING
	newOrder := models.Order{
		CustomerName: input.CustomerName,
		Amount:       input.Amount,
		Status:       "PENDING",
	}

	// 3. store new order to database
	if err := h.DB.Create(&newOrder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create order"})
		return
	}

	// 4. push new order ID to JobChan for worker to process
	h.JobChan <- newOrder.ID

	// 5. respond to client immediately that order is received and being processed
	c.JSON(http.StatusAccepted, gin.H{
		"message":  "order received and being processed",
		"order_id": newOrder.ID,
		"status":   newOrder.Status,
	})
}

// GetOrders list all orders in database
func (h *OrderHandler) GetOrders(c *gin.Context) {
	var orders []models.Order
	h.DB.Find(&orders)
	c.JSON(http.StatusOK, orders)
}
