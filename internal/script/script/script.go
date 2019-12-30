package script

import (
	"time"
)

const (
	DefaultInterval time.Duration = time.Second * 30
)

type Script struct {
	Name     string
	Body     []byte
	Interval time.Duration
}
