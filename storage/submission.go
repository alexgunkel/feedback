package storage

import "gorm.io/gorm"

type Submission struct {
	gorm.Model
	SurveyID uint
	Ratings  []Rating
	Texts    []Text
}

type Rating struct {
	gorm.Model
	SubmissionID uint
	Category     string
	Value        int
}

type Text struct {
	gorm.Model
	SubmissionID uint
	Category     string
	Value        string `sql:"type:LONGTEXT CHARACTER SET utf8 COLLATE utf8_general_ci"`
}
