package models

import (
	"gorm.io/gorm"
)

type PR struct {
	gorm.Model
	Title     string
	Author    string
	Reviewer  string
	NextActor string
}
