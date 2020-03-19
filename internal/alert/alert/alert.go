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

var (
	alertsPool = sync.Pool{}
)

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

func ReleaseAlert(a *Alert) {
	a.reset()
	alertsPool.Put(a)
}

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

func (l Level) String() string {
	switch l {
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

func (a *Alert) SetName(name string) {
	a.name = name
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

func (a *Alert) Name() string {
	return a.name
}
