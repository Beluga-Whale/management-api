package models

import (
	"time"

	"gorm.io/gorm"
)

type Status string

const (
	Active Status ="active"
	Inactive Status = "inactive"
)

type Priority string

const (
	Low Priority = "low"
	Medium Priority = "medium"
	High Priority = "high"
)

type Tasks struct {
	gorm.Model
	DueDate time.Time
	Title string
	Description string
	Status  Status `gorm:"type:task_status;not null; default:'active'"`
	Completed bool
	Priority Priority `gorm:"type:task_priority;not null;default:'low'"`
	UserID uint //NOTE - FK
}