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

// GetOrder get order by ID
func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.Order
	if err := h.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy đơn hàng"})
		return
	}
	c.JSON(http.StatusOK, order)
}

// UpdateOrder update order by ID (not implemented in this example, but you can add if needed)
func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.Order
	if err := h.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Đơn hàng không tồn tại"})
		return
	}

	var input struct {
		CustomerName string  `json:"customer_name"`
		Amount       float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.DB.Model(&order).Updates(input)
	c.JSON(http.StatusOK, order)
}

// DeleteOrder soft delete order by ID (not implemented in this example, but you can add if needed)
func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	id := c.Param("id")
	if err := h.DB.Delete(&models.Order{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi xóa đơn hàng"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Đã xóa đơn hàng #" + id})
}
