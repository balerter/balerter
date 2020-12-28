package service

type Service struct {
	Address string `json:"address" yaml:"address"`
}

func (s *Service) Validate() error {
	return nil
}
