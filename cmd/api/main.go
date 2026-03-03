package main

import (
	"log"
	"logistics-simulator/internal/database"
	"logistics-simulator/internal/handlers"
	"logistics-simulator/internal/workers"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Khởi tạo DB
	database.InitDB()

	// 2. Tạo Channel cho Worker Pool (Buffer 100)
	jobChan := make(chan uint, 100)

	// 3. Chạy Worker Pool (3 thợ)
	workers.StartWorkerPool(3, jobChan, database.DB)

	// 4. Khởi tạo Handler
	orderHandler := &handlers.OrderHandler{
		DB:      database.DB,
		JobChan: jobChan,
	}

	// 5. Setup Gin Router
	router := gin.Default()

	router.POST("/orders", orderHandler.CreateOrder)
	router.GET("/orders", orderHandler.GetOrders)
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 6. Chạy Server ở port 8080
	log.Println("is running")
	router.Run()
}
