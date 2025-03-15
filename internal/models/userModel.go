package models

import "gorm.io/gorm"

type Role string

const (
	Admin Role = "admin"
	User  Role = "user"
)

type Users struct {
	gorm.Model
	Email string `gorm:"unique"`
	Name string  `gorm:"not null"`
	Password string `gorm: "unique;not null"`
	Photo string `gorm:"default:'https://images.unsplash.com/photo-1438761681033-6461ffad8d80?q=80&w=2070&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D'"`
	Bio string 
	Role  Role `gorm:"type:user_role;not null;default:'user'"`
	isVerified bool 
	Tasks []Tasks `gorm:"foreignKey:UserID"`
}