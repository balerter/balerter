package alert

import (
	"sync"
	"time"
)

type Level int

const (
	LevelSuccess Level = 1
	LevelError   Level = 2
	LevelWarn    Level = 3
)

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
