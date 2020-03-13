package memory

import "fmt"

func (m *Memory) Put(name string, val string) error {
	m.mxKV.Lock()
	defer m.mxKV.Unlock()

	if _, ok := m.kv[name]; ok {
		return fmt.Errorf("variable already exists")
	}

	m.kv[name] = val

	return nil
}

func (m *Memory) Upsert(name string, val string) error {
	m.mxKV.Lock()
	defer m.mxKV.Unlock()

	m.kv[name] = val

	return nil
}

func (m *Memory) Get(name string) (string, error) {
	m.mxKV.RLock()
	defer m.mxKV.RUnlock()

	v, ok := m.kv[name]
	if !ok {
		return "", fmt.Errorf("variable not exists")
	}

	return v, nil
}

func (m *Memory) Delete(name string) error {
	m.mxKV.Lock()
	defer m.mxKV.Unlock()

	if _, ok := m.kv[name]; !ok {
		return fmt.Errorf("variable not exists")
	}

	delete(m.kv, name)

	return nil
}
