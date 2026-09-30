package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"device-fleet-monitor/internal/device"
)

type Handler struct {
	store *device.Store
}

func NewHandler(store *device.Store) *Handler {
	return &Handler{
		store: store,
	}
}

type registerDeviceRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type deviceResponse struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Status        string    `json:"status"`
	LastHeartbeat time.Time `json:"last_heartbeat,omitempty"`
}

func (h *Handler) RegisterDevice(w http.ResponseWriter, r *http.Request) {
	var req registerDeviceRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	req.ID = strings.TrimSpace(req.ID)
	req.Name = strings.TrimSpace(req.Name)

	if req.ID == "" || req.Name == "" {
		http.Error(w, "id and name are required", http.StatusBadRequest)
		return
	}

	d := device.Device{
		ID:   req.ID,
		Name: req.Name,
	}

	if err := h.store.Add(d); err != nil {
		if err == device.ErrDeviceExists {
			http.Error(w, "device already exists", http.StatusConflict)
			return
		}

		http.Error(w, "failed to register device", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, toDeviceResponse(d, time.Now()))
}

func (h *Handler) Heartbeat(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if id == "" {
		http.Error(w, "device id is required", http.StatusBadRequest)
		return
	}

	now := time.Now()

	if err := h.store.UpdateHeartbeat(id, now); err != nil {
		if err == device.ErrDeviceNotFound {
			http.Error(w, "device not found", http.StatusNotFound)
			return
		}

		http.Error(w, "failed to update heartbeat", http.StatusInternalServerError)
		return
	}

	d, err := h.store.Get(id)
	if err != nil {
		http.Error(w, "device not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, toDeviceResponse(d, now))
}

func (h *Handler) GetDevice(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	d, err := h.store.Get(id)
	if err != nil {
		if err == device.ErrDeviceNotFound {
			http.Error(w, "device not found", http.StatusNotFound)
			return
		}

		http.Error(w, "failed to get device", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, toDeviceResponse(d, time.Now()))
}

func (h *Handler) ListDevices(w http.ResponseWriter, r *http.Request) {
	devices := h.store.List()
	now := time.Now()

	response := make([]deviceResponse, 0, len(devices))

	for _, d := range devices {
		response = append(response, toDeviceResponse(d, now))
	}

	writeJSON(w, http.StatusOK, response)
}

func (h *Handler) Summary(w http.ResponseWriter, r *http.Request) {
	devices := h.store.List()
	now := time.Now()

	summary := struct {
		Total   int `json:"total"`
		Online  int `json:"online"`
		Offline int `json:"offline"`
	}{
		Total: len(devices),
	}

	for _, d := range devices {
		if d.Status(now) == "ONLINE" {
			summary.Online++
		} else {
			summary.Offline++
		}
	}

	writeJSON(w, http.StatusOK, summary)
}

func toDeviceResponse(d device.Device, now time.Time) deviceResponse {
	return deviceResponse{
		ID:            d.ID,
		Name:          d.Name,
		Status:        d.Status(now),
		LastHeartbeat: d.LastHeartbeat,
	}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		return
	}
}
