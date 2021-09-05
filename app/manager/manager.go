package manager

import "github.com/arsmn/ontest/settings"

type Manager struct {
	stg settings.Provider
}

func NewManager(s settings.Provider) *Manager {
	m := &Manager{}

	m.stg = s

	return m
}
