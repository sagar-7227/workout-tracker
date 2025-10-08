package controllers

import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReportController struct{ DB *gorm.DB }

// @Summary Workout report summary
// @Tags Reports
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]any
// @Router /reports [get]
func (rc *ReportController) Summary(c *gin.Context) {
	uid := c.GetString("userID")
	var total int64
	rc.DB.Table("workouts").Where("user_id = ?", uid).Count(&total)

	var lastWeek int64
	oneWeekAgo := time.Now().Add(-7 * 24 * time.Hour)
	rc.DB.Table("workouts").Where("user_id = ? AND created_at >= ?", uid, oneWeekAgo).Count(&lastWeek)

	// Most common muscle group in scheduled week (simple heuristic via exercises.category)
	var top struct {
		Category string
		Cnt      int64
	}
	rc.DB.Raw(`SELECT e.category, COUNT(*) AS cnt
FROM workout_exercises we
JOIN workouts w ON w.id = we.workout_id
JOIN exercises e ON e.id = we.exercise_id
WHERE w.user_id = ?
GROUP BY e.category
ORDER BY cnt DESC
LIMIT 1`, uid).Scan(&top)

	c.JSON(http.StatusOK, gin.H{
		"total_workouts":       total,
		"last_week_sessions":   lastWeek,
		"most_common_category": top.Category,
	})
}
