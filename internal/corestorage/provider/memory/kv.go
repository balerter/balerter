package memory

import "fmt"

func (m *storageKV) Put(name, val string) error {
	m.mxKV.Lock()
	defer m.mxKV.Unlock()

	if _, ok := m.kv[name]; ok {
		return fmt.Errorf("variable already exists")
	}

	m.kv[name] = val

	return nil
}

func (m *storageKV) Upsert(name, val string) error {
	m.mxKV.Lock()
	defer m.mxKV.Unlock()

	m.kv[name] = val

	return nil
}

func (m *storageKV) Get(name string) (string, error) {
	m.mxKV.RLock()
	defer m.mxKV.RUnlock()

	v, ok := m.kv[name]
	if !ok {
		return "", fmt.Errorf("variable not exists")
	}

	return v, nil
}

func (m *storageKV) Delete(name string) error {
	m.mxKV.Lock()
	defer m.mxKV.Unlock()

	if _, ok := m.kv[name]; !ok {
		return fmt.Errorf("variable not exists")
	}

	delete(m.kv, name)

	return nil
}

func (m *storageKV) All() (map[string]string, error) {
	m.mxKV.RLock()
	defer m.mxKV.RUnlock()

	r := make(map[string]string)

	for name, value := range m.kv {
		r[name] = value
	}

	return r, nil
}
