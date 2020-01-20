package manager

type AlertInfo struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func (m *Manager) GetAlerts() []*AlertInfo {
	m.activeMx.RLock()
	defer m.activeMx.RUnlock()

	info := make([]*AlertInfo, 0)

	for name, count := range m.active {
		info = append(info, &AlertInfo{
			Name:  name,
			Count: count,
		})
	}

	return info
}
