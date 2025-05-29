package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PR struct {
	gorm.Model
	Title     string
	Author    string
	Reviewers datatypes.JSON `json:"reviewers"`
	NextActor string
}
