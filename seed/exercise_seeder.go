package seed

import (
	"log"
	"gorm.io/gorm"
	"workout-tracker/models"
)

func RunExerciseSeeder(db *gorm.DB) {
	var count int64
	db.Model(&models.Exercise{}).Count(&count)
	if count > 0 {
		return
	}

	items := []models.Exercise{
		{Name: "Squat", Description: "Barbell back squat", Category: "legs"},
		{Name: "Bench Press", Description: "Barbell bench press", Category: "chest"},
		{Name: "Deadlift", Description: "Conventional deadlift", Category: "back"},
		{Name: "Overhead Press", Description: "Standing OHP", Category: "shoulders"},
		{Name: "Pull-Up", Description: "Bodyweight pull-up", Category: "back"},
		{Name: "Running", Description: "Treadmill running", Category: "cardio"},
		{Name: "Plank", Description: "Core plank", Category: "core"},
	}
	if err := db.Create(&items).Error; err != nil {
		log.Printf("seed error: %v", err)
	}
}
