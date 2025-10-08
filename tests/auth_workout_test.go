package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"workout-tracker/models"
	"workout-tracker/routes"
)

type loginResp struct {
	Token string `json:"token"`
}

func setupTestRouter() (*gin.Engine, *gorm.DB) {
	_ = os.Setenv("JWT_SECRET", "test-secret")
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.User{}, &models.Exercise{}, &models.Workout{}, &models.WorkoutExercise{})
	// seed minimal
	db.Create(&models.Exercise{Name: "Squat", Category: "legs"})

	r := gin.Default()
	routes.Register(r, db)
	return r, db
}

func TestSignupLoginAndCreateWorkout(t *testing.T) {
	r, _ := setupTestRouter()

	// signup
	req := httptest.NewRequest("POST", "/api/auth/signup", strings.NewReader(`{"name":"A","email":"a@a.com","password":"secret1"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("signup: %d", rec.Code)
	}

	// login
	req = httptest.NewRequest("POST", "/api/auth/login", strings.NewReader(`{"email":"a@a.com","password":"secret1"}`))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("login: %d", rec.Code)
	}
	var lr loginResp
	_ = json.Unmarshal(rec.Body.Bytes(), &lr)
	if lr.Token == "" {
		t.Fatal("empty token")
	}

	// create workout
	payload := `{"title":"Leg Day","exercises":[{"exercise_id":"1","sets":3,"repetitions":10,"weight":60}]}`
	req = httptest.NewRequest("POST", "/api/workouts", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+lr.Token)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create workout: %d %s", rec.Code, rec.Body.String())
	}
}
