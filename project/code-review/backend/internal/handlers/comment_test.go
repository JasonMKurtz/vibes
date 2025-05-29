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

func setupCommentDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	if err := db.AutoMigrate(&models.PR{}, &models.Comment{}); err != nil {
		t.Fatalf("failed to migrate db: %v", err)
	}
	return db
}

func setupCommentRouter(h *CommentHandler) *gin.Engine {
	router := gin.Default()
	router.GET("/prs/:id/comments", h.ListComments)
	router.POST("/prs/:id/comments", h.CreateComment)
	return router
}

func TestCreateAndListComments(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupCommentDB(t)

	pr := models.PR{Title: "Test", Author: "alice"}
	if err := db.Create(&pr).Error; err != nil {
		t.Fatalf("create pr: %v", err)
	}

	h := NewCommentHandler(db)
	r := setupCommentRouter(h)

	comment := models.Comment{Author: "bob", Body: "looks good"}
	body, _ := json.Marshal(comment)
	url := "/prs/" + strconv.Itoa(int(pr.ID)) + "/comments"
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
	var resp []models.Comment
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(resp) != 1 || resp[0].Body != comment.Body {
		t.Fatalf("unexpected response: %#v", resp)
	}
}
