package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"code-review/internal/models"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	if err := db.AutoMigrate(&models.PR{}, &models.Review{}, &models.Comment{}); err != nil {
		t.Fatalf("failed to migrate db: %v", err)
	}
	return db
}

func setupRouter(h *PRHandler) *gin.Engine {
	router := gin.Default()
	router.GET("/prs", h.ListPRs)
	router.POST("/prs", h.CreatePR)
	router.PUT("/prs/:id/next", h.UpdateNextActor)
	return router
}

func toJSONSlice(t *testing.T, vals []string) datatypes.JSON {
	t.Helper()
	b, err := json.Marshal(vals)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return datatypes.JSON(b)
}

func TestCreateAndListPRs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cases := []struct {
		name string
		pr   models.PR
	}{
		{
			name: "single reviewer",
			pr:   models.PR{Title: "Test", Author: "alice", Reviewers: toJSONSlice(t, []string{"bob"})},
		},
		{
			name: "multiple reviewers",
			pr:   models.PR{Title: "Another", Author: "alice", Reviewers: toJSONSlice(t, []string{"bob", "carol"})},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db := setupTestDB(t)
			h := NewPRHandler(db)
			r := setupRouter(h)

			body, _ := json.Marshal(tc.pr)
			req := httptest.NewRequest(http.MethodPost, "/prs", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			if w.Code != http.StatusCreated {
				t.Fatalf("expected %d got %d", http.StatusCreated, w.Code)
			}

			req = httptest.NewRequest(http.MethodGet, "/prs", nil)
			w = httptest.NewRecorder()
			r.ServeHTTP(w, req)
			if w.Code != http.StatusOK {
				t.Fatalf("expected %d got %d", http.StatusOK, w.Code)
			}
			var resp []models.PR
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if len(resp) != 1 {
				t.Fatalf("expected 1 PR got %d", len(resp))
			}
			if resp[0].Title != tc.pr.Title {
				t.Fatalf("expected title %q got %q", tc.pr.Title, resp[0].Title)
			}
		})
	}
}

func TestUpdateNextActor(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cases := []struct {
		name       string
		nextActor  string
		statusCode int
	}{
		{"author valid", "alice", http.StatusOK},
		{"reviewer valid", "bob", http.StatusOK},
		{"invalid", "mallory", http.StatusBadRequest},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db := setupTestDB(t)
			reviewers := toJSONSlice(t, []string{"bob", "carol"})
			basePR := models.PR{Title: "Update", Author: "alice", Reviewers: reviewers, NextActor: "alice"}
			if err := db.Create(&basePR).Error; err != nil {
				t.Fatalf("create pr: %v", err)
			}
			h := NewPRHandler(db)
			r := setupRouter(h)

			payload := map[string]string{"next_actor": tc.nextActor}
			body, _ := json.Marshal(payload)
			url := "/prs/" + strconv.Itoa(int(basePR.ID)) + "/next"
			req := httptest.NewRequest(http.MethodPut, url, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			if w.Code != tc.statusCode {
				t.Fatalf("expected %d got %d", tc.statusCode, w.Code)
			}
			var pr models.PR
			if err := db.First(&pr, basePR.ID).Error; err != nil {
				t.Fatalf("fetch pr: %v", err)
			}
			if tc.statusCode == http.StatusOK && pr.NextActor != tc.nextActor {
				t.Fatalf("expected next_actor %q got %q", tc.nextActor, pr.NextActor)
			}
		})
	}
}
