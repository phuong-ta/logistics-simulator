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

	// 6. Chạy Server ở port 8080
	log.Println("Server đang chạy tại port 8080...")
	router.Run(":8080")
}
