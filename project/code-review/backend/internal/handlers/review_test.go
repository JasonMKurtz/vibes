package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"code-review/internal/models"
)

func setupReviewDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	if err := db.AutoMigrate(&models.PR{}, &models.Review{}); err != nil {
		t.Fatalf("failed to migrate db: %v", err)
	}
	return db
}

func setupReviewRouter(h *ReviewHandler) *gin.Engine {
	router := gin.Default()
	router.GET("/prs/:id/reviews", h.ListReviews)
	router.POST("/prs/:id/reviews", h.CreateReview)
	return router
}

func TestCreateAndListReviews(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupReviewDB(t)

	pr := models.PR{Title: "Test", Author: "alice"}
	if err := db.Create(&pr).Error; err != nil {
		t.Fatalf("create pr: %v", err)
	}

	h := NewReviewHandler(db)
	r := setupReviewRouter(h)

	review := models.Review{Reviewer: "bob", State: "approved"}
	body, _ := json.Marshal(review)
	url := "/prs/" + strconv.Itoa(int(pr.ID)) + "/reviews"
	req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected %d got %d", http.StatusCreated, w.Code)
	}

	req = httptest.NewRequest(http.MethodGet, url, nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected %d got %d", http.StatusOK, w.Code)
	}
	var resp []models.Review
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(resp) != 1 || resp[0].State != review.State || resp[0].Reviewer != review.Reviewer {
		t.Fatalf("unexpected response: %#v", resp)
	}
}
