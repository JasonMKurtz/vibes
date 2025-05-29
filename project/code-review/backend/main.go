package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"code-review/internal/handlers"
	"code-review/internal/models"
)

func setupDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("code_review.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&models.PR{}); err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	db, err := setupDB()
	if err != nil {
		log.Fatalf("failed to setup db: %v", err)
	}

	router := gin.Default()

	prHandler := handlers.NewPRHandler(db)
	router.GET("/prs", prHandler.ListPRs)
	router.POST("/prs", prHandler.CreatePR)
	router.PUT("/prs/:id/next", prHandler.UpdateNextActor)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
