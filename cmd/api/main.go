package main

import (
	"net/http"

	"log"
	"logistics-simulator/internal/database"
	"logistics-simulator/internal/workers"

	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDB()
	log.Println("--- Kết nối thành công! App đã sẵn sàng cho bước tiếp theo ---")
	// 2. create a channel to send order IDs to workers ( 100 buffer)
	jobChan := make(chan uint, 100)

	// 3. Khởi tạo Worker Pool với 5 thợ (Goroutines)
	workers.StartWorkerPool(5, jobChan, database.DB)

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.Run()
}
