package models

import "time"

type Workout struct {
	ID          string            `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID      string            `gorm:"type:uuid;index" json:"user_id"`
	Title       string            `gorm:"size:100" json:"title"`
	Notes       string            `gorm:"type:text" json:"notes"`
	ScheduledAt *time.Time        `json:"scheduled_at"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Exercises   []WorkoutExercise `json:"exercises"`
}
