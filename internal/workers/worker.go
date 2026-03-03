package workers

import (
	"log"
	"logistics-simulator/internal/models"
	"time"

	"gorm.io/gorm"
)

// ProcessOrder simulate the processing of an order by a worker
func ProcessOrder(workerID int, orderID uint, db *gorm.DB) {
	log.Printf("[Worker %d] start to process order #%d\n", workerID, orderID)

	// 1. update status to PROCESSING
	db.Model(&models.Order{}).Where("id = ?", orderID).Update("status", "PROCESSING")

	// 2. simulate time to process order (paking, shipping, etc.)
	// set 10 seconds to see the effect clearly in logs, you can adjust as needed
	time.Sleep(10 * time.Second)

	// 3. update status to COMPLETED
	db.Model(&models.Order{}).Where("id = ?", orderID).Update("status", "COMPLETED")

	log.Printf("[Worker %d] completed order #%d ✅\n", workerID, orderID)
}

// StartWorkerPool create worker (Goroutines) to listen on jobChan and process orders concurrently
func StartWorkerPool(numWorkers int, jobChan <-chan uint, db *gorm.DB) {
	for i := 1; i <= numWorkers; i++ {
		go func(id int) {
			log.Printf("[Worker %d] ready to process orders!\n", id)

			// Worker alway listen channel jobChan
			for orderID := range jobChan {
				ProcessOrder(id, orderID, db)
			}
		}(i)
	}
}
