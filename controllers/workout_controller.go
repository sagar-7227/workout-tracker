package controllers

import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"workout-tracker/models"
)

type WorkoutController struct{ DB *gorm.DB }

type WorkoutExerciseDTO struct {
	ExerciseID  string  `json:"exercise_id" binding:"required"`
	Sets        int     `json:"sets" binding:"required,min=1"`
	Repetitions int     `json:"repetitions" binding:"required,min=1"`
	Weight      float64 `json:"weight"`
}

type CreateWorkoutReq struct {
	Title       string               `json:"title" binding:"required,min=2"`
	Notes       string               `json:"notes"`
	ScheduledAt *time.Time           `json:"scheduled_at"`
	Exercises   []WorkoutExerciseDTO `json:"exercises" binding:"required,min=1,dive"`
}

// @Summary Create workout
// @Tags Workouts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body CreateWorkoutReq true "Workout"
// @Success 201 {object} models.Workout
// @Router /workouts [post]
func (wc *WorkoutController) Create(c *gin.Context) {
	uid := c.GetString("userID")
	var req CreateWorkoutReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	w := models.Workout{UserID: uid, Title: req.Title, Notes: req.Notes, ScheduledAt: req.ScheduledAt}
	for _, ex := range req.Exercises {
		w.Exercises = append(w.Exercises, models.WorkoutExercise{ExerciseID: ex.ExerciseID, Sets: ex.Sets, Reps: ex.Repetitions, Weight: ex.Weight})
	}
	if err := wc.DB.Create(&w).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, w)
}

// @Summary List workouts (sorted by scheduled_at asc, nulls last)
// @Tags Workouts
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.Workout
// @Router /workouts [get]
func (wc *WorkoutController) List(c *gin.Context) {
	uid := c.GetString("userID")
	var ws []models.Workout
	wc.DB.Preload("Exercises").Where("user_id = ?", uid).Order("scheduled_at NULLS LAST, created_at desc").Find(&ws)
	c.JSON(http.StatusOK, ws)
}

// @Summary Get workout by ID
// @Tags Workouts
// @Security BearerAuth
// @Produce json
// @Param id path string true "Workout ID"
// @Success 200 {object} models.Workout
// @Router /workouts/{id} [get]
func (wc *WorkoutController) Get(c *gin.Context) {
	uid := c.GetString("userID")
	id := c.Param("id")
	var w models.Workout
	if err := wc.DB.Preload("Exercises").Where("id = ? AND user_id = ?", id, uid).First(&w).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, w)
}

type UpdateWorkoutReq struct {
	Title       *string               `json:"title"`
	Notes       *string               `json:"notes"`
	ScheduledAt **time.Time           `json:"scheduled_at"`
	Exercises   *[]WorkoutExerciseDTO `json:"exercises"`
}

// @Summary Update workout
// @Tags Workouts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Workout ID"
// @Param payload body UpdateWorkoutReq true "Update"
// @Success 200 {object} models.Workout
// @Router /workouts/{id} [put]
func (wc *WorkoutController) Update(c *gin.Context) {
	uid := c.GetString("userID")
	id := c.Param("id")
	var w models.Workout
	if err := wc.DB.Where("id = ? AND user_id = ?", id, uid).First(&w).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	var req UpdateWorkoutReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Title != nil {
		w.Title = *req.Title
	}
	if req.Notes != nil {
		w.Notes = *req.Notes
	}
	if req.ScheduledAt != nil {
		w.ScheduledAt = *req.ScheduledAt
	}
	if req.Exercises != nil {
		// Replace all exercises (simple approach)
		wc.DB.Where("workout_id = ?", w.ID).Delete(&models.WorkoutExercise{})
		for _, ex := range *req.Exercises {
			w.Exercises = append(w.Exercises, models.WorkoutExercise{ExerciseID: ex.ExerciseID, Sets: ex.Sets, Reps: ex.Repetitions, Weight: ex.Weight})
		}
	}
	if err := wc.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(&w).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, w)
}

// @Summary Delete workout
// @Tags Workouts
// @Security BearerAuth
// @Param id path string true "Workout ID"
// @Success 204
// @Router /workouts/{id} [delete]
func (wc *WorkoutController) Delete(c *gin.Context) {
	uid := c.GetString("userID")
	id := c.Param("id")
	res := wc.DB.Where("id = ? AND user_id = ?", id, uid).Delete(&models.Workout{})
	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
