package alert

import (
	"sync"
	"time"
)

type Level int

const (
	LevelInfo    Level = 0
	LevelSuccess Level = 1
	LevelError   Level = 2
	LevelWarn    Level = 3
)

type Alert struct {
	mx sync.RWMutex

	Level      Level
	LastChange time.Time
	Start      time.Time
}

func New() *Alert {
	now := time.Now()

	a := &Alert{
		Level:      LevelInfo,
		LastChange: now,
		Start:      now,
	}

	return a
}

func (a *Alert) UpdateLevel(level Level) {
	a.mx.Lock()
	defer a.mx.RUnlock()

	a.Level = level
	a.LastChange = time.Now()
}
