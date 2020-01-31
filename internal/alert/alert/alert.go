package alert

import "time"

type Level int

const (
	LevelInfo    Level = 0
	LevelSuccess Level = 1
	LevelError   Level = 2
	LevelWarn    Level = 3
)

type Alert struct {
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
