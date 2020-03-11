package alert

import (
	"fmt"
	"sync"
	"time"
)

type Level int

const (
	LevelSuccess Level = 1
	LevelWarn    Level = 2
	LevelError   Level = 3
)

func LevelFromString(s string) (Level, error) {
	switch s {
	case "success":
		return LevelSuccess, nil
	case "warn", "warning":
		return LevelWarn, nil
	case "error":
		return LevelError, nil
	}

	return 0, fmt.Errorf("bad level")
}

func (l *Level) String() string {
	switch *l {
	case LevelSuccess:
		return "success"
	case LevelWarn:
		return "warning"
	case LevelError:
		return "error"
	}

	panic("unexpected level value")
}

type Alert struct {
	mx sync.RWMutex

	level      Level
	lastChange time.Time
	start      time.Time
	count      int
}

func New() *Alert {
	now := time.Now()

	a := &Alert{
		level:      LevelSuccess,
		lastChange: now,
		start:      now,
		count:      0,
	}

	return a
}

func (a *Alert) UpdateLevel(level Level) {
	a.mx.Lock()
	defer a.mx.Unlock()

	a.level = level
	a.lastChange = time.Now()
	a.count = 0
}

func (a *Alert) Inc() {
	a.mx.Lock()
	defer a.mx.Unlock()

	a.count++
}

func (a *Alert) Level() Level {
	a.mx.Lock()
	defer a.mx.Unlock()

	l := a.level

	return l
}

func (a *Alert) Count() int {
	a.mx.Lock()
	defer a.mx.Unlock()

	c := a.count

	return c
}

func (a *Alert) GetLastChangeTime() time.Time {
	a.mx.Lock()
	defer a.mx.Unlock()

	r := a.lastChange

	return r
}

func (a *Alert) GetStartTime() time.Time {
	a.mx.Lock()
	defer a.mx.Unlock()

	r := a.start

	return r
}
