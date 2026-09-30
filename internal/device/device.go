package device

import "time"

const heartbeatTimeout = 30 * time.Second

type Device struct {
	ID            string
	Name          string
	LastHeartbeat time.Time
}

func (d Device) Status(now time.Time) string {
	if d.LastHeartbeat.IsZero() {
		return "OFFLINE"
	}

	if now.Sub(d.LastHeartbeat) <= heartbeatTimeout {
		return "ONLINE"
	}

	return "OFFLINE"
}
