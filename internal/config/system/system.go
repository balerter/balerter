package system

type System struct {
	JobWorkersCount int `json:"jobWorkersCount" yaml:"jobWorkersCount" hcl:"jobWorkersCount,optional"`
}

func (s *System) Validate() error {
	return nil
}
