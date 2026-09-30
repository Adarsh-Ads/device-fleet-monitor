package main

import (
	"fmt"
	"net/http"

	"device-fleet-monitor/internal/api"
	"device-fleet-monitor/internal/device"
)

func main() {
	store := device.NewStore()
	handler := api.NewHandler(store)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /devices", handler.RegisterDevice)
	mux.HandleFunc("POST /devices/{id}/heartbeat", handler.Heartbeat)
	mux.HandleFunc("GET /devices", handler.ListDevices)
	mux.HandleFunc("GET /devices/{id}", handler.GetDevice)
	mux.HandleFunc("GET /summary", handler.Summary)

	fmt.Println("Starting Device Fleet Monitor on :8080")

	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		fmt.Println("Server failed:", err)
	}
}
