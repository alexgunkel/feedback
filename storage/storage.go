package storage

import (
	"fmt"
	"github.com/alexgunkel/feedback/config"
	"github.com/alexgunkel/feedback/util"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Input struct {
	Val    url.Values
	Survey uint
}

type Storage struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewStorage(cfg *config.Config) *Storage {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User, cfg.Database.Password, "127.0.0.1", 3306, cfg.Database.Database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	util.PanicOnError(err)
	util.PanicOnError(db.AutoMigrate(&Survey{}, &Submission{}, &Rating{}, &Text{}))
	util.PanicOnError(err)

	return &Storage{db: db, cfg: cfg}
}

func (s *Storage) HandleRatings(c <-chan Input) {
	for v := range c {
		newSubmission := &Submission{SurveyID: v.Survey}
		for _, chapter := range s.cfg.Chapters {
			for _, question := range chapter.Questions {
				if value, ok := v.Val[question.Short]; ok && len(value) == 1 {
					if question.Type == config.TextField {
						s.addTextToSubmission(newSubmission, question, value[0])
					} else {
						s.addRatingToSubmission(value[0], newSubmission, question.Short)
					}
				}
			}
		}

		s.db.Save(newSubmission)
	}
}

func (s *Storage) addTextToSubmission(newSubmission *Submission, question config.Question, value string) {
	if len(value) < 1 {
		return
	}
	newSubmission.Texts = append(newSubmission.Texts, Text{Category: question.Short, Value: value})
}

func (s *Storage) addRatingToSubmission(value string, newSubmission *Submission, key string) {
	if i, err := strconv.Atoi(value); err == nil {
		newSubmission.Ratings = append(newSubmission.Ratings, Rating{Category: key, Value: i})
	}
}

type StatisticsResult struct {
	ID       uint
	Title    string
	DateTime time.Time
	Category string
	Value    int
	Cnt      int
}

func (s *Storage) GetStatistics() []Submission {
	var result []Submission
	tx := s.db.Preload("Ratings").Preload("Texts").Find(&result)

	if tx.Error != nil {
		panic(tx.Error)
	}

	return result
}

type Result struct {
	Title   string
	Ratings map[string]map[int]int
	Texts   map[string][]string
}

func HtmlEscape(in string) string {
	return "<p>" + strings.Replace(in, "\n", "</p><p>", -1) + "</p>"
}

func (s *Storage) GetResult(id Survey) (result Result) {
	result.Title = id.Title
	ratingStatistics := make([]struct {
		Category string
		Value    int
		Cnt      int
	}, 0)
	s.db.Raw("SELECT r.category, r.value, COUNT(s.id) as cnt FROM submissions s INNER JOIN feedback.ratings r on s.id = r.submission_id WHERE s.survey_id = ? GROUP BY r.category, r.value;", id.ID).Find(&ratingStatistics)

	result.Ratings = map[string]map[int]int{}
	result.Texts = map[string][]string{}
	for _, next := range ratingStatistics {
		if _, ok := result.Ratings[next.Category]; !ok {
			result.Ratings[next.Category] = map[int]int{}
		}
		result.Ratings[next.Category][next.Value] = next.Cnt
	}

	textStatistics := make([]struct {
		Category string
		Value    string
	}, 0)
	s.db.Raw("SELECT t.category, t.value FROM submissions s INNER JOIN feedback.texts t on s.id = t.submission_id WHERE s.survey_id = ? GROUP BY t.category, t.value;", id.ID).Find(&textStatistics)

	for _, next := range textStatistics {
		if _, ok := result.Texts[next.Category]; !ok {
			result.Texts[next.Category] = make([]string, 0)
		}
		result.Texts[next.Category] = append(result.Texts[next.Category], next.Value)
	}

	return result
}

func (s *Storage) GetSurveys() (res []Survey) {
	ret := s.db.Preload("Submissions").Find(&res)
	util.PanicOnError(ret.Error)
	return res
}

func (s *Storage) GetSurvey(id int) *Survey {
	survey := new(Survey)
	if err := s.db.First(survey, id).Error; err != nil {
		return nil
	}

	return survey
}

func (s *Storage) Add(survey *Survey) {
	s.db.Save(survey)
}
