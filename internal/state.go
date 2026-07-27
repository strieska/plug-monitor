package internal

import (
	"sync"
	"time"
)

type MonitorState struct {
	Running   bool
	StartTime time.Time
	Notified  bool
}

type Monitor struct {
	mu    sync.Mutex
	State MonitorState
}

func NewMonitor() *Monitor {
	return &Monitor{
		State: MonitorState{
			Running:  false,
			Notified: false,
		},
	}
}
