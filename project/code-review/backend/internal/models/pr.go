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
	Files     datatypes.JSON `json:"files"`
}

// FileDiff describes the before and after contents of a file in a PR.
type FileDiff struct {
	Filename string `json:"filename"`
	Before   string `json:"before"`
	After    string `json:"after"`
}
