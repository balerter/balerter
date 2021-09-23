package system

import (
	"fmt"
	"time"
)

type System struct {
	JobWorkersCount int    `json:"jobWorkersCount" yaml:"jobWorkersCount" hcl:"jobWorkersCount,optional"`
	CronLocation    string `json:"cronLocation" yaml:"cronLocation" hcl:"cronLocation,optional"`
}

func (s *System) Validate() error {
	if s.CronLocation != "" {
		_, err := time.LoadLocation(s.CronLocation)
		if err != nil {
			return fmt.Errorf("error parse cronLocation, %w", err)
		}
	}
	return nil
}
