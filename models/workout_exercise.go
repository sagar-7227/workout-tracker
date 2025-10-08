package models

type WorkoutExercise struct {
	ID         string  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	WorkoutID  string  `gorm:"type:uuid;index" json:"workout_id"`
	ExerciseID string  `gorm:"type:uuid;index" json:"exercise_id"`
	Sets       int     `json:"sets"`
	Reps       int     `json:"repetitions"`
	Weight     float64 `json:"weight"`
}
