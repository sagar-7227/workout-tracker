package config

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"workout-tracker/models"
)

func GetPort() string {
	p := os.Getenv("APP_PORT")
	if p == "" {
		p = "8080"
	}
	return p
}

func MustInitDB() *gorm.DB {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		panic("DB_DSN not provided")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Exercise{},
		&models.Workout{},
		&models.WorkoutExercise{},
	); err != nil {
		panic(fmt.Errorf("automigrate: %w", err))
	}

	return db
}
