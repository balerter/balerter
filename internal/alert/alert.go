package alert

import (
	"errors"
	"time"
)

type Options struct {
	Fields   []string
	Channels []string
	Quiet    bool
	Repeat   int
	Image    string
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

// Alert is base struct for store Alert information
type Alert struct {
	//mx sync.RWMutex

	Name       string
	Level      Level
	LastChange time.Time
	Start      time.Time
	Count      int
}

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

// UpdateLevel allows to update an Alert Level
//func (a *Alert) UpdateLevel(Level Level) {
//	a.mx.Lock()
//	defer a.mx.Unlock()
//
//	a.Level = Level
//	a.LastChange = time.Now()
//	a.Count = 0
//}
//
//// Inc increments an Alert counter
//func (a *Alert) Inc() {
//	a.mx.Lock()
//	defer a.mx.Unlock()
//
//	a.Count++
//}
//
//// HasLevel return true if the alert.Level equals to the Level argument
//func (a *Alert) HasLevel(Level Level) bool {
//	a.mx.Lock()
//	defer a.mx.Unlock()
//
//	return a.Level == Level
//}
//
//// Level allows to get an alert Level
//func (a *Alert) Level() Level {
//	a.mx.Lock()
//	defer a.mx.Unlock()
//
//	l := a.Level
//
//	return l
//}
//
//// Count allows to get an alert counter value
//func (a *Alert) Count() int {
//	a.mx.Lock()
//	defer a.mx.Unlock()
//
//	c := a.Count
//
//	return c
//}
//
//// GetLastChangeTime allows to get an Alert.LastChange value
//func (a *Alert) GetLastChangeTime() time.Time {
//	a.mx.Lock()
//	defer a.mx.Unlock()
//
//	r := a.LastChange
//
//	return r
//}
//
//// GetStartTime allows to get an Alert.Start value
//func (a *Alert) GetStartTime() time.Time {
//	a.mx.Lock()
//	defer a.mx.Unlock()
//
//	r := a.Start
//
//	return r
//}
//
//// Name returns an Alert Name
//func (a *Alert) Name() string {
//	return a.Name
//}
