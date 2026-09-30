package device

import (
	"testing"
	"time"
)

func TestDeviceStatus(t *testing.T) {
	now := time.Date(2026, 9, 30, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		lastHeartbeat time.Time
		expected      string
	}{
		{
			name:     "no heartbeat",
			expected: "OFFLINE",
		},
		{
			name:          "heartbeat 29 seconds ago",
			lastHeartbeat: now.Add(-29 * time.Second),
			expected:      "ONLINE",
		},
		{
			name:          "heartbeat exactly 30 seconds ago",
			lastHeartbeat: now.Add(-30 * time.Second),
			expected:      "ONLINE",
		},
		{
			name:          "heartbeat 31 seconds ago",
			lastHeartbeat: now.Add(-31 * time.Second),
			expected:      "OFFLINE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Device{
				ID:            "device-01",
				Name:          "Test Device",
				LastHeartbeat: tt.lastHeartbeat,
			}

			got := d.Status(now)

			if got != tt.expected {
				t.Errorf("expected status %q, got %q", tt.expected, got)
			}
		})
	}
}
