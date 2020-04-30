package script

import (
	"crypto/sha1"
	"fmt"

	"github.com/robfig/cron/v3"
)

var DefaultSchedule *Schedule

func init() {
	var err error
	DefaultSchedule, err = NewSchedule("@every 60s")
	if err != nil {
		panic(err)
	}
}

func New() *Script {
	return &Script{}
}

type Schedule struct {
	cron.Schedule
	spec string
}

func NewSchedule(spec string) (*Schedule, error) {
	sc, err := cron.ParseStandard(spec)
	if err != nil {
		return nil, err
	}

	return &Schedule{
		Schedule: sc,
		spec:     spec,
	}, nil
}

func (sc *Schedule) String() string {
	return sc.spec
}

type Script struct {
	Name       string
	Body       []byte
	Schedule   *Schedule
	Ignore     bool
	Channels   []string
	IsTest     bool
	TestTarget string
}

func (s *Script) Hash() string {
	return fmt.Sprintf("%x", sha1.Sum(append([]byte(s.Name+"@"), s.Body...)))
}
