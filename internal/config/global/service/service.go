package service

type Service struct {
	Address string `json:"address" yaml:"address"`
	Metrics bool   `json:"metrics" yaml:"metrics"`
}

func (s *Service) Validate() error {
	return nil
}
