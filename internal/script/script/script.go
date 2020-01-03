package script

import (
	"crypto/sha1"
	"fmt"
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

func (s *Script) Hash() string {
	return fmt.Sprintf("%x", sha1.Sum(append([]byte(s.Name+"@"), s.Body...)))
}
