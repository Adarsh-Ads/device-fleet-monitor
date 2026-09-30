package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"time"
)

type simulatedDevice struct {
	ID   string
	Name string
}

var devices = []simulatedDevice{
	{ID: "device-01", Name: "Temperature Sensor"},
	{ID: "device-02", Name: "Pressure Sensor"},
	{ID: "device-03", Name: "Humidity Sensor"},
	{ID: "device-04", Name: "Motion Sensor"},
	{ID: "device-05", Name: "Power Monitor"},
}

func registerDevice(client *http.Client, baseURL string, d simulatedDevice) error {
	body := map[string]string{
		"id":   d.ID,
		"name": d.Name,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	resp, err := client.Post(
		baseURL+"/devices",
		"application/json",
		bytes.NewReader(data),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusConflict {
		fmt.Printf("%s is already registered\n", d.ID)
		return nil
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf(
			"registration failed for %s: status %d",
			d.ID,
			resp.StatusCode,
		)
	}

	fmt.Printf("Registered %s\n", d.ID)

	return nil
}

func sendHeartbeat(
	client *http.Client,
	baseURL string,
	id string,
) error {
	resp, err := client.Post(
		baseURL+"/devices/"+id+"/heartbeat",
		"",
		nil,
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"heartbeat failed for %s: status %d",
			id,
			resp.StatusCode,
		)
	}

	fmt.Printf("Heartbeat sent: %s\n", id)

	return nil
}

func runDevice(
	client *http.Client,
	baseURL string,
	d simulatedDevice,
	stop <-chan struct{},
) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	if err := sendHeartbeat(client, baseURL, d.ID); err != nil {
		fmt.Printf("Heartbeat error for %s: %v\n", d.ID, err)
	}

	for {
		select {
		case <-ticker.C:
			if err := sendHeartbeat(client, baseURL, d.ID); err != nil {
				fmt.Printf(
					"Heartbeat error for %s: %v\n",
					d.ID,
					err,
				)
			}

		case <-stop:
			fmt.Printf("Stopped simulator for %s\n", d.ID)
			return
		}
	}
}

func main() {
	stopDevice := flag.String(
		"stop",
		"device-05",
		"device whose heartbeat simulator should stop",
	)

	stopAfter := flag.Duration(
		"stop-after",
		40*time.Second,
		"time after which the selected device stops sending heartbeats",
	)

	flag.Parse()

	baseURL := "http://localhost:8080"

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	for _, d := range devices {
		if err := registerDevice(client, baseURL, d); err != nil {
			fmt.Printf(
				"Could not register %s: %v\n",
				d.ID,
				err,
			)
			return
		}
	}

	stops := make(map[string]chan struct{})

	for _, d := range devices {
		stops[d.ID] = make(chan struct{})

		go runDevice(
			client,
			baseURL,
			d,
			stops[d.ID],
		)
	}

	fmt.Printf(
		"All simulators started. %s will stop after %s.\n",
		*stopDevice,
		*stopAfter,
	)

	time.Sleep(*stopAfter)

	stop, exists := stops[*stopDevice]
	if !exists {
		fmt.Printf("Unknown device: %s\n", *stopDevice)
		return
	}

	fmt.Printf(
		"Stopping heartbeat simulator for %s...\n",
		*stopDevice,
	)

	close(stop)

	fmt.Println(
		"The device should become OFFLINE after 30 seconds without a heartbeat.",
	)

	select {}
}
