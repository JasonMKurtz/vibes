package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"code-review/internal/models"
)

type PRHandler struct {
	DB *gorm.DB
}

func NewPRHandler(db *gorm.DB) *PRHandler {
	return &PRHandler{DB: db}
}

func (h *PRHandler) ListPRs(c *gin.Context) {
	var prs []models.PR
	if err := h.DB.Find(&prs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, prs)
}

func (h *PRHandler) CreatePR(c *gin.Context) {
	var pr models.PR
	if err := c.ShouldBindJSON(&pr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.DB.Create(&pr).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, pr)
}

func (h *PRHandler) UpdateNextActor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var payload struct {
		NextActor string `json:"next_actor"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var pr models.PR
	if err := h.DB.First(&pr, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "PR not found"})
		return
	}
	pr.NextActor = payload.NextActor
	if err := h.DB.Save(&pr).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pr)
}
