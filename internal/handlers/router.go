package handlers

import "github.com/gin-gonic/gin"

// SetupRouter and return Gin Engine with all routes setup
func SetupRouter(h *OrderHandler) *gin.Engine {
	r := gin.Default()

	// Group all order-related routes under /orders
	orderRoutes := r.Group("/orders")
	{
		orderRoutes.POST("", h.CreateOrder)       // Create
		orderRoutes.GET("", h.GetOrders)          // Read All
		orderRoutes.GET("/:id", h.GetOrder)       // Read One
		orderRoutes.PUT("/:id", h.UpdateOrder)    // Update
		orderRoutes.DELETE("/:id", h.DeleteOrder) // Delete
	}

	return r
}
