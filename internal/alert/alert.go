package alert

import (
	"errors"
	"strconv"
	"time"
)

// Options represents Alert options
type Options struct {
	Channels []string          `json:"channels"`
	Quiet    bool              `json:"quiet"`
	Repeat   int               `json:"repeat"`
	Image    string            `json:"image"`
	Fields   map[string]string `json:"fields"`
	Group    string            `json:"group"`
}

func NewOptions() *Options {
	return &Options{
		Fields: map[string]string{},
	}
}

// Level is the type for describe an Alert Level
type Level int

const (
	// LevelSuccess is Success Level of an alert
	LevelSuccess Level = 1
	// LevelWarn is Waring Level of an alert
	LevelWarn Level = 2
	// LevelError is Error Level of an alert
	LevelError Level = 3
)

var (
	// ErrBadLevel represent an error if user provide the incorrect level value
	ErrBadLevel = errors.New("bad level")
)

const (
	levelStringSuccess  = "success"
	levelStringWarning1 = "warning"
	levelStringWarning2 = "warn"
	levelStringError    = "error"
)

// LevelFromString returns Level based on provided string or error
func LevelFromString(s string) (Level, error) {
	switch s {
	case levelStringSuccess:
		return LevelSuccess, nil
	case levelStringWarning1, levelStringWarning2:
		return LevelWarn, nil
	case levelStringError:
		return LevelError, nil
	}

	return 0, ErrBadLevel
}

// LevelFromInt returns Level based on provided int or error
func LevelFromInt(i int) (Level, error) {
	switch i {
	case int(LevelSuccess):
		return LevelSuccess, nil
	case int(LevelWarn):
		return LevelWarn, nil
	case int(LevelError):
		return LevelError, nil
	}

	return 0, ErrBadLevel
}

// NumString returns numeric value of the Level as a string
func (l Level) NumString() string {
	return strconv.Itoa(int(l))
}

// String returns string value of the Level
func (l Level) String() string {
	switch l {
	case LevelSuccess:
		return levelStringSuccess
	case LevelWarn:
		return levelStringWarning1
	case LevelError:
		return levelStringError
	}

	panic("unexpected Level value")
}

// Alerts contains slice of alerts
type Alerts []*Alert

// Alert is base struct for store Alert information
type Alert struct {
	Name       string    `json:"name"`
	Level      Level     `json:"level"`
	LastChange time.Time `json:"last_change"`
	Start      time.Time `json:"start"`
	Count      int       `json:"count"`
}

// New creates new Alert
func New(name string) *Alert {
	now := time.Now()

	a := &Alert{
		Name:       name,
		Level:      LevelSuccess,
		LastChange: now,
		Start:      now,
		Count:      0,
	}

	return a
}
