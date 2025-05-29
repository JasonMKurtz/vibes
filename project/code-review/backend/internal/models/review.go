package models

import "gorm.io/gorm"

// Review represents a review on a PR. It can hold an overall state
// like approval or request for changes.
type Review struct {
	gorm.Model
	PRID     uint
	Reviewer string
	State    string
	Comments []Comment `gorm:"constraint:OnDelete:CASCADE"`
}

// Comment is a single comment left on a PR or as part of a review.
type Comment struct {
	gorm.Model
	PRID     uint
	ReviewID *uint
	File     string
	Line     *int
	Author   string
	Body     string
}
