package config

import (
	"github.com/alexgunkel/feedback/util"
	"gopkg.in/yaml.v2"
	"os"
)

type QuestionType string

const (
	TextField QuestionType = "text"
)

type Question struct {
	Title string       `yaml:"title"`
	Short string       `yaml:"short"`
	Type  QuestionType `yaml:"type"`
}

func (q Question) IsText() bool {
	return q.Type == "text"
}

type Chapter struct {
	Title     string     `yaml:"title"`
	Questions []Question `yaml:"questions"`
}

type Option struct {
	Description string `yaml:"description"`
	ID          int    `yaml:"id"`
}

type DbConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type Config struct {
	Listen   string    `yaml:"listen"`
	Domain   string    `yaml:"domain"`
	Chapters []Chapter `yaml:"chapters"`
	Options  []Option  `yaml:"options"`
	Database DbConfig  `yaml:"db_config"`
}

func ReadConfig() *Config {
	file := "config.yaml"
	if len(os.Args) > 1 {
		file = os.Args[1]
	}
	cfg, err := os.ReadFile(file)
	util.PanicOnError(err)
	form := new(Config)
	err = yaml.Unmarshal(cfg, form)
	util.PanicOnError(err)
	return form
}
