# Mini Device Fleet Monitor

A small Go application that monitors a fleet of simulated devices through periodic heartbeats.

Each device sends a heartbeat to the server. The server records the latest heartbeat and dynamically determines whether the device is `ONLINE` or `OFFLINE`.

## Features

- Device registration
- Heartbeat ingestion
- Dynamic ONLINE/OFFLINE status
- Fleet listing
- Individual device lookup
- Fleet summary
- Concurrent device simulator
- Automated tests
- Concurrency-safe in-memory storage
- Configurable simulator stop behavior

## Architecture

```text
                    ┌─────────────────────┐
                    │   Device Simulator  │
                    │                     │
                    │  5 simulated devices│
                    └──────────┬──────────┘
                               │
                         HTTP heartbeats
                               │
                               ▼
                    ┌─────────────────────┐
                    │    HTTP API Server  │
                    │      net/http       │
                    └──────────┬──────────┘
                               │
                               ▼
                    ┌─────────────────────┐
                    │     Device Store    │
                    │  map + RWMutex      │
                    └─────────────────────┘