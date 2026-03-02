package main

import (
	"net/http"

	"log"
	"logistics-simulator/internal/database"

	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDB()
	log.Println("--- Kết nối thành công! App đã sẵn sàng cho bước tiếp theo ---")
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.Run()
}
