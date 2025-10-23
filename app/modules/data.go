package data

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Структура для конфига
type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
}

type Resume struct {
	Main           MainInfo        `yaml:"main"`
	Skills         Skills          `yaml:"skills"`
	Tools          []string        `yaml:"tools"`
	Experience     []Experience    `yaml:"experience"`
	Education      []Education     `yaml:"education"`
	Certifications []Certification `yaml:"certifications"`
	Projects       []Project       `yaml:"projects"`
	About          []About         `yaml:"about"`
}
type Skills struct {
	ProgrammingLanguages string `yaml:"Programming Languages"`
	Frameworks           string `yaml:"Frameworks"`
	Tools                string `yaml:"Tools"`
}
type MainInfo struct {
	Name     string `yaml:"name"`
	Email    string `yaml:"email"`
	GitHub   string `yaml:"github"`
	LinkedIn string `yaml:"linkedin"`
	JobTitle string `yaml:"job_title"`
}

type Experience struct {
	Title       string `yaml:"title"`
	Company     string `yaml:"company"`
	Location    string `yaml:"location"`
	StartDate   string `yaml:"start_date"`
	EndDate     string `yaml:"end_date"`
	Description string `yaml:"description"`
}

type Education struct {
	Degree     string `yaml:"degree"`
	University string `yaml:"university"`
	Location   string `yaml:"location"`
	StartDate  string `yaml:"start_date"`
	EndDate    string `yaml:"end_date"`
}

type Certification struct {
	Name         string `yaml:"name"`
	Organization string `yaml:"organization"`
	Date         string `yaml:"date"`
	Description  string `yaml:"description"`
}

type Project struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Link        string `yaml:"link,omitempty"`
}

type About struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
}

// Функция загрузки
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func LoadData(path string) (*Resume, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	var resume Resume
	if err := yaml.Unmarshal(data, &resume); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	return &resume, nil
}
