package handlers

import (
	"encoding/json"
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

// ShowDiff returns the diff for all files in a PR.
func (h *PRHandler) ShowDiff(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var pr models.PR
	if err := h.DB.First(&pr, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "PR not found"})
		return
	}

	var files []models.FileDiff
	if len(pr.Files) > 0 {
		if err := json.Unmarshal(pr.Files, &files); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid files"})
			return
		}
	}

	type fileResp struct {
		Filename string          `json:"filename"`
		Diff     []models.DiffOp `json:"diff"`
	}
	var resp []fileResp
	for _, f := range files {
		diff := models.DiffLines(f.Before, f.After)
		resp = append(resp, fileResp{Filename: f.Filename, Diff: diff})
	}
	c.JSON(http.StatusOK, resp)
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

	var reviewers []string
	if err := json.Unmarshal(pr.Reviewers, &reviewers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid reviewers"})
		return
	}

	valid := payload.NextActor == pr.Author
	if !valid {
		for _, r := range reviewers {
			if payload.NextActor == r {
				valid = true
				break
			}
		}
	}
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid next_actor"})
		return
	}

	pr.NextActor = payload.NextActor
	if err := h.DB.Save(&pr).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pr)
}
