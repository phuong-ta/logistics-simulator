package main

import (
	"log"
	"logistics-simulator/internal/database"
	"logistics-simulator/internal/handlers"
	"logistics-simulator/internal/workers"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1.init database connection and auto-migrate Order model
	database.InitDB()

	// 2. Create Channel for Worker Pool (Buffer 100)
	jobChan := make(chan uint, 100)

	// 3. Run Worker Pool (3 workers)
	workers.StartWorkerPool(3, jobChan, database.DB)

	// 4. init Handler
	orderHandler := &handlers.OrderHandler{
		DB:      database.DB,
		JobChan: jobChan,
	}

	// 5. Setup Gin Router
	router := gin.Default()

	router.POST("/orders", orderHandler.CreateOrder)
	router.GET("/orders", orderHandler.GetOrders)

	// 6. Chạy Server ở port 8080
	log.Println("Server is running on port 8080...")
	router.Run(":8080")
}
