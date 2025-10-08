package models

import "time"

type Exercise struct {
	ID          string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name        string    `gorm:"size:100;uniqueIndex" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Category    string    `gorm:"size:50" json:"category"` // muscle group or category
	CreatedAt   time.Time `json:"created_at"`
}
