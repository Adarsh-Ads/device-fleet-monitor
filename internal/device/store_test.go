package device

import (
	"testing"
	"time"
)

func TestStoreAddAndGet(t *testing.T) {
	store := NewStore()

	d := Device{
		ID:   "device-01",
		Name: "Test Device",
	}

	err := store.Add(d)
	if err != nil {
		t.Fatalf("unexpected error adding device: %v", err)
	}

	got, err := store.Get("device-01")
	if err != nil {
		t.Fatalf("unexpected error getting device: %v", err)
	}

	if got.ID != d.ID {
		t.Errorf("expected ID %q, got %q", d.ID, got.ID)
	}

	if got.Name != d.Name {
		t.Errorf("expected name %q, got %q", d.Name, got.Name)
	}
}

func TestStoreRejectsDuplicateDevice(t *testing.T) {
	store := NewStore()

	d := Device{
		ID:   "device-01",
		Name: "Test Device",
	}

	if err := store.Add(d); err != nil {
		t.Fatalf("unexpected error adding first device: %v", err)
	}

	err := store.Add(d)

	if err != ErrDeviceExists {
		t.Errorf("expected ErrDeviceExists, got %v", err)
	}
}

func TestStoreGetUnknownDevice(t *testing.T) {
	store := NewStore()

	_, err := store.Get("does-not-exist")

	if err != ErrDeviceNotFound {
		t.Errorf("expected ErrDeviceNotFound, got %v", err)
	}
}

func TestStoreUpdateHeartbeat(t *testing.T) {
	store := NewStore()

	d := Device{
		ID:   "device-01",
		Name: "Test Device",
	}

	if err := store.Add(d); err != nil {
		t.Fatalf("unexpected error adding device: %v", err)
	}

	heartbeat := time.Date(
		2026,
		9,
		30,
		12,
		0,
		0,
		0,
		time.UTC,
	)

	if err := store.UpdateHeartbeat("device-01", heartbeat); err != nil {
		t.Fatalf("unexpected error updating heartbeat: %v", err)
	}

	got, err := store.Get("device-01")
	if err != nil {
		t.Fatalf("unexpected error getting device: %v", err)
	}

	if !got.LastHeartbeat.Equal(heartbeat) {
		t.Errorf(
			"expected heartbeat %v, got %v",
			heartbeat,
			got.LastHeartbeat,
		)
	}
}
