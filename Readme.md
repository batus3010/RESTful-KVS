# Key-Value-Store

A simple RESTful key-value store implemented in Go, containerized with Docker, backed by a file-system store, built following TDD approach with a test suite for unit and integration tests. Future enhancements include clustering for fault tolerance, replication, and consensus via Raft.

---

## Table of Contents

* [Features](#features)
* [Requirements](#requirements)
* [Architecture](#architecture)
* [Getting Started](#getting-started)

    * [Clone the Repository](#clone-the-repository)
    * [Build](#build)
    * [Run Locally](#run-locally)
    * [Run with Docker](#run-with-docker)
* [Usage](#usage)

    * [API Endpoints](#api-endpoints)
* [Testing](#testing)
* [Project Structure](#project-structure)
* [Future Work](#future-work)
* [License](#license)

---

## Features

* **RESTful API** for basic key-value operations (`GET`, `PUT`, `DELETE`).
* **File-system-backed storage** with simple persistence.
* **Unit and integration tests** for core components and server.
* **Dockerized** for easy deployment.

---

## Requirements

* Go 1.18+
* Docker (optional, for containerized deployment)

---

## Architecture

The codebase is organized into:

* **`KVS.go`**: Core interface and implementation of the in-memory key-value store.
* **`file_system_store.go`**: File-system-backed implementation of the KVS interface.
* **`server.go`**, **`cmd/RESTful/main.go`**: HTTP server exposing RESTful endpoints.
* **`Dockerfile`**: Defines the container image.
* **Tests**: Unit tests (`*_test.go`) and an integration test (`KVS_integration_test.go`).

---

## Getting Started

### Clone the Repository

```bash
git clone https://github.com/yourusername/key-value-store.git
cd key-value-store
```

### Build

```bash
go build ./cmd/RESTful
```

### Run Locally

```bash
# Start server on default port 5000
go run cmd/RESTful/main.go
```

### Run with Docker

```bash
# Build Docker image
docker build -t kvs-server .

# Run container, mapping port 5000
docker run -p 5000:5000 --name kvs-server kvs-server
```

---

## Usage

### API Endpoints

| Method | Path        | Description                             |
|--------|-------------|-----------------------------------------|
| GET    | `/kv/{key}` | Retrieve the value for `key`.           |
| PUT    | `/kv/{key}` | Set the value (in request body)         |
| DELETE | `/kv/{key}` | Delete the specified `key`.             |
| ALL    | `/all`      | Get all key-value pair in JSON response |

**Example**:

```bash
# Set key "foo"
curl -X PUT localhost:5000/kv/foo -d 'bar'

# Get key "foo"
curl localhost:5000/kv/foo

# Delete key "foo"
curl -X DELETE localhost:5000/kv/foo
```

---

## Testing

Run all tests:

```bash
go test ./...
```

This covers:

* Unit tests for `KVS` and file-system store.
* Integration test for the RESTful server.

---

## Project Structure

```text
Key-Value-Store/
├─ cmd/
│  └─ RESTful/
│     └─ main.go         # Entry point for HTTP server
├─ KVS.go               # Core key-value store interface & in-memory impl
├─ file_system_store.go # File-system-backed store
├─ server.go            # HTTP handlers wiring
├─ Dockerfile           # Container image definition
├─ *_test.go            # Unit & integration tests
└─ README.md            # Project documentation
```

---

## Future Work

0. **In progress**: Persist the data in disk, probably by a Json file
1. **Fault Tolerance & Replication**: Integrate a consensus algorithm (e.g., Raft) to replicate state across nodes and achieve leader election.
2. **Distributed Deployment**: Container orchestration with Kubernetes StatefulSets, service discovery, and sharding.
3. **Tunable Consistency**: Expose consistency levels (ONE, QUORUM, ALL) and implement anti-entropy for eventual convergence.
4. **Monitoring & Metrics**: Add Prometheus metrics, structured logging, and health probes.

---

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
