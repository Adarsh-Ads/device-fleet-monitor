package device

import (
	"errors"
	"sync"
	"time"
)

var ErrDeviceNotFound = errors.New("device not found")
var ErrDeviceExists = errors.New("device already exists")

type Store struct {
	mu      sync.RWMutex
	devices map[string]Device
}

func NewStore() *Store {
	return &Store{
		devices: make(map[string]Device),
	}
}

func (s *Store) Add(device Device) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.devices[device.ID]; exists {
		return ErrDeviceExists
	}

	s.devices[device.ID] = device
	return nil
}

func (s *Store) Get(id string) (Device, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	device, exists := s.devices[id]
	if !exists {
		return Device{}, ErrDeviceNotFound
	}

	return device, nil
}

func (s *Store) UpdateHeartbeat(id string, heartbeatTime time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	device, exists := s.devices[id]
	if !exists {
		return ErrDeviceNotFound
	}

	device.LastHeartbeat = heartbeatTime
	s.devices[id] = device

	return nil
}

func (s *Store) List() []Device {
	s.mu.RLock()
	defer s.mu.RUnlock()

	devices := make([]Device, 0, len(s.devices))

	for _, device := range s.devices {
		devices = append(devices, device)
	}

	return devices
}
