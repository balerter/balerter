package memory

import "github.com/balerter/balerter/internal/alert"

func (m *storageAlert) Update(name string, level alert.Level) (*alert.Alert, bool, error) {
	panic("memory provider not implemented")
}

//func (m *storageAlert) GetOrNew(name string) (*alert.Alert, error) {
//	m.mxAlerts.RLock()
//	a, ok := m.alerts[name]
//	m.mxAlerts.RUnlock()
//	if ok {
//		return a, nil
//	}
//
//	m.mxAlerts.Lock()
//	defer m.mxAlerts.Unlock()
//
//	a, ok = m.alerts[name]
//	if ok {
//		return a, nil
//	}
//	a = alert.New(name)
//	m.alerts[name] = a
//	return a, nil
//}
//
//func (m *storageAlert) All() ([]*alert.Alert, error) {
//	var result []*alert.Alert
//
//	m.mxAlerts.RLock()
//	defer m.mxAlerts.RUnlock()
//
//	for _, a := range m.alerts {
//		result = append(result, a)
//	}
//
//	return result, nil
//}
//
//func (m *storageAlert) Store(a *alert.Alert) error {
//	m.mxAlerts.Lock()
//	m.alerts[a.Name()] = a
//	m.mxAlerts.Unlock()
//	return nil
//}
//
//func (m *storageAlert) Get(name string) (*alert.Alert, error) {
//	m.mxAlerts.RLock()
//	defer m.mxAlerts.RUnlock()
//
//	a, ok := m.alerts[name]
//	if !ok {
//		return nil, fmt.Errorf("alert not found")
//	}
//
//	return a, nil
//}
