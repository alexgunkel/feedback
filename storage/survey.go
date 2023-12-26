package storage

import (
	"fmt"
	"github.com/goodsign/monday"
	"gorm.io/gorm"
	"time"
)

type Survey struct {
	gorm.Model
	Title         string
	AccessKey     string
	EvaluationKey string
	DateTime      time.Time
	Submissions   []Submission
}

func (s Survey) GetAccessPath() string {
	return fmt.Sprintf("/%d/%s", s.ID, s.AccessKey)
}
func (s Survey) GetEvaluationPath() string {
	return fmt.Sprintf("/%d/%s", s.ID, s.EvaluationKey)
}

const DateTimeFormat = "Monday, den 02. January 2006, um 15:04"

func formatDateTime(in time.Time) string {
	return monday.Format(in, DateTimeFormat, monday.LocaleDeDE)
}

func (s Survey) GetDateTime() string {
	return formatDateTime(s.DateTime) //s.DateTime.Format(DateTimeFormat)
}

func (s Survey) GetCreatedAt() string {
	return formatDateTime(s.CreatedAt)
}
