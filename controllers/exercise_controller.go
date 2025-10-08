package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ExerciseController struct{ DB *gorm.DB }

// @Summary List exercises
// @Tags Exercises
// @Security BearerAuth
// @Produce json
// @Success 200 {array} map[string]any
// @Router /exercises [get]
func (e *ExerciseController) List(c *gin.Context) {
	var items []map[string]any
	e.DB.Model(&struct{}{}).Table("exercises").Find(&items)
	c.JSON(http.StatusOK, items)
}
