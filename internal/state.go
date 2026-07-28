package internal

import (
	"sync"
	"time"
)

type MonitorState struct {
	Running   bool      `json:"running"`
	StartTime time.Time `json:"start_time"`
	Notified  bool      `json:"notified"`
}

type Monitor struct {
	mu    sync.Mutex
	State MonitorState
}

func NewMonitor() *Monitor {
	return &Monitor{}
}

// Update updates the monitor state.
//
// Returns:
//   - current state
//   - should notify
//   - running duration
func (m *Monitor) Update(
	power float64,
	threshold time.Duration,
) (MonitorState, bool, time.Duration) {

	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()

	if power <= 0 {

		m.State.Running = false
		m.State.StartTime = time.Time{}
		m.State.Notified = false

		return m.State, false, 0
	}

	if !m.State.Running {

		m.State.Running = true
		m.State.StartTime = now
		m.State.Notified = false

		return m.State, false, 0
	}

	running := now.Sub(m.State.StartTime)

	if !m.State.Notified &&
		running >= threshold {

		m.State.Notified = true

		return m.State, true, running
	}

	return m.State, false, running
}

func (m *Monitor) GetState() MonitorState {

	m.mu.Lock()
	defer m.mu.Unlock()

	return m.State
}
