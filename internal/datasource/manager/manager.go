package manager

type Manager struct {
}

func New() (*Manager, error) {
	m := &Manager{}

	return m, nil
}
