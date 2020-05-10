package alert

import (
	"fmt"
	"sync"
	"time"
)

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
	alertsPool = sync.Pool{}
)

// AcquireAlert returns new Alert from a sync.Pool
func AcquireAlert() *Alert {
	a := alertsPool.Get()
	if a == nil {
		a := &Alert{}
		a.init()
		return a
	}

	a.(*Alert).init()

	return a.(*Alert)
}

// ReleaseAlert puts an alert to a sync.pool
func ReleaseAlert(a *Alert) {
	a.reset()
	alertsPool.Put(a)
}

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

	return 0, fmt.Errorf("bad level")
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

func (a *Alert) init() {
	now := time.Now()

	a.level = LevelSuccess
	a.lastChange = now
	a.start = now
	a.count = 0
}

func (a *Alert) reset() {
	a.name = ""
	a.level = 0
	a.lastChange = time.Time{}
	a.start = time.Time{}
	a.count = 0
}

// SetName allows to set an Alert name
func (a *Alert) SetName(name string) {
	a.name = name
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
