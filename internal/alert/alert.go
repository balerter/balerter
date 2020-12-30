package alert

import (
	"errors"
	"sync"
	"time"
)

type Options struct {
	Fields   []string
	Channels []string
	Quiet    bool
	Repeat   int
	Image    string
}

// Level is the type for describe an Alert level
type Level int

const (
	// LevelSuccess is Success level of an alert
	LevelSuccess Level = 1
	// LevelWarn is Waring level of an alert
	LevelWarn Level = 2
	// LevelError is Error level of an alert
	LevelError Level = 3
)

var (
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

	panic("unexpected level value")
}

// Alert is base struct for store Alert information
type Alert struct {
	mx sync.RWMutex

	name       string
	level      Level
	lastChange time.Time
	start      time.Time
	count      int
}

func New(name string) *Alert {
	now := time.Now()

	a := &Alert{
		name:       name,
		level:      LevelSuccess,
		lastChange: now,
		start:      now,
	}

	return a
}

// UpdateLevel allows to update an Alert level
func (a *Alert) UpdateLevel(level Level) {
	a.mx.Lock()
	defer a.mx.Unlock()

	a.level = level
	a.lastChange = time.Now()
	a.count = 0
}

// Inc increments an Alert counter
func (a *Alert) Inc() {
	a.mx.Lock()
	defer a.mx.Unlock()

	a.count++
}

// HasLevel return true if the alert.Level equals to the level argument
func (a *Alert) HasLevel(level Level) bool {
	a.mx.Lock()
	defer a.mx.Unlock()

	return a.level == level
}

// Level allows to get an alert level
func (a *Alert) Level() Level {
	a.mx.Lock()
	defer a.mx.Unlock()

	l := a.level

	return l
}

// Count allows to get an alert counter value
func (a *Alert) Count() int {
	a.mx.Lock()
	defer a.mx.Unlock()

	c := a.count

	return c
}

// GetLastChangeTime allows to get an Alert.LastChange value
func (a *Alert) GetLastChangeTime() time.Time {
	a.mx.Lock()
	defer a.mx.Unlock()

	r := a.lastChange

	return r
}

// GetStartTime allows to get an Alert.Start value
func (a *Alert) GetStartTime() time.Time {
	a.mx.Lock()
	defer a.mx.Unlock()

	r := a.start

	return r
}

// Name returns an Alert name
func (a *Alert) Name() string {
	return a.name
}
