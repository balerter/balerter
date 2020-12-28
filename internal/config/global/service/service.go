package service

type Service struct {
	Address string
}

func (s *Service) Validate() error {
	return nil
}
