package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"device-fleet-monitor/internal/device"
)

func newTestServer() *http.ServeMux {
	store := device.NewStore()
	handler := NewHandler(store)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /devices", handler.RegisterDevice)
	mux.HandleFunc("POST /devices/{id}/heartbeat", handler.Heartbeat)
	mux.HandleFunc("GET /devices", handler.ListDevices)
	mux.HandleFunc("GET /devices/{id}", handler.GetDevice)
	mux.HandleFunc("GET /summary", handler.Summary)

	return mux
}

func TestRegisterDevice(t *testing.T) {
	mux := newTestServer()

	body := `{"id":"device-01","name":"Temperature Sensor"}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/devices",
		strings.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusCreated,
			rec.Code,
		)
	}

	var response deviceResponse

	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.ID != "device-01" {
		t.Errorf("expected ID device-01, got %q", response.ID)
	}

	if response.Name != "Temperature Sensor" {
		t.Errorf(
			"expected name Temperature Sensor, got %q",
			response.Name,
		)
	}

	if response.Status != "OFFLINE" {
		t.Errorf(
			"expected status OFFLINE, got %q",
			response.Status,
		)
	}
}

func TestHeartbeat(t *testing.T) {
	mux := newTestServer()

	registerBody := `{"id":"device-01","name":"Temperature Sensor"}`

	registerReq := httptest.NewRequest(
		http.MethodPost,
		"/devices",
		strings.NewReader(registerBody),
	)

	registerReq.Header.Set("Content-Type", "application/json")

	registerRec := httptest.NewRecorder()

	mux.ServeHTTP(registerRec, registerReq)

	if registerRec.Code != http.StatusCreated {
		t.Fatalf(
			"expected registration status %d, got %d",
			http.StatusCreated,
			registerRec.Code,
		)
	}

	heartbeatReq := httptest.NewRequest(
		http.MethodPost,
		"/devices/device-01/heartbeat",
		nil,
	)

	heartbeatRec := httptest.NewRecorder()

	mux.ServeHTTP(heartbeatRec, heartbeatReq)

	if heartbeatRec.Code != http.StatusOK {
		t.Fatalf(
			"expected heartbeat status %d, got %d",
			http.StatusOK,
			heartbeatRec.Code,
		)
	}

	var response deviceResponse

	if err := json.NewDecoder(heartbeatRec.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Status != "ONLINE" {
		t.Errorf(
			"expected status ONLINE, got %q",
			response.Status,
		)
	}
}

func TestHeartbeatUnknownDevice(t *testing.T) {
	mux := newTestServer()

	req := httptest.NewRequest(
		http.MethodPost,
		"/devices/unknown/heartbeat",
		nil,
	)

	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusNotFound,
			rec.Code,
		)
	}
}

func TestDuplicateRegistration(t *testing.T) {
	mux := newTestServer()

	body := `{"id":"device-01","name":"Temperature Sensor"}`

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(
			http.MethodPost,
			"/devices",
			strings.NewReader(body),
		)

		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if i == 0 && rec.Code != http.StatusCreated {
			t.Fatalf(
				"first registration: expected %d, got %d",
				http.StatusCreated,
				rec.Code,
			)
		}

		if i == 1 && rec.Code != http.StatusConflict {
			t.Fatalf(
				"duplicate registration: expected %d, got %d",
				http.StatusConflict,
				rec.Code,
			)
		}
	}
}

func TestSummary(t *testing.T) {
	mux := newTestServer()

	devices := []string{
		`{"id":"device-01","name":"Device 1"}`,
		`{"id":"device-02","name":"Device 2"}`,
		`{"id":"device-03","name":"Device 3"}`,
	}

	for _, body := range devices {
		req := httptest.NewRequest(
			http.MethodPost,
			"/devices",
			strings.NewReader(body),
		)

		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf(
				"expected registration status %d, got %d",
				http.StatusCreated,
				rec.Code,
			)
		}
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/devices/device-01/heartbeat",
		nil,
	)

	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected heartbeat status %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}

	summaryReq := httptest.NewRequest(
		http.MethodGet,
		"/summary",
		nil,
	)

	summaryRec := httptest.NewRecorder()

	mux.ServeHTTP(summaryRec, summaryReq)

	if summaryRec.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			summaryRec.Code,
		)
	}

	var summary struct {
		Total   int `json:"total"`
		Online  int `json:"online"`
		Offline int `json:"offline"`
	}

	if err := json.NewDecoder(summaryRec.Body).Decode(&summary); err != nil {
		t.Fatalf("failed to decode summary: %v", err)
	}

	if summary.Total != 3 {
		t.Errorf("expected total 3, got %d", summary.Total)
	}

	if summary.Online != 1 {
		t.Errorf("expected online 1, got %d", summary.Online)
	}

	if summary.Offline != 2 {
		t.Errorf("expected offline 2, got %d", summary.Offline)
	}
}
