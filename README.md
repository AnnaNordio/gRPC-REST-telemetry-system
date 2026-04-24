# Telemetry Benchmarking Suite: REST vs. gRPC

This project was developed for the **Runtimes for Concurrency and Distribution** course at the **University of Padua**, Academic Year **2025/2026**.

The suite is a high-precision tool designed to evaluate and compare the performance characteristics of **REST (JSON over HTTP/1.1)** and **gRPC (Protobuf over HTTP/2)** in high-frequency telemetry scenarios. It allows for real-time monitoring of latency, throughput, and protocol overhead.

---

## Prerequisites

Ensure you have the following tools installed:

### 1. GNU Make
Used to automate builds and environment management.
* **Ubuntu/Debian:** `sudo apt update && sudo apt install build-essential`
* **MacOS:** `brew install make`

### 2. Docker & Docker Compose
The system runs in a containerized environment to ensure network isolation and reproducibility.
* **Installation:** [Get Docker](https://docs.docker.com/get-docker/)
* **Verification:** `docker compose version`

---

## Getting Started

Follow these steps to initialize and run the suite.

### 1. Build and Setup
Generate dependencies and build the Docker images:
```bash
make docker-gen
```
### 2. Interactive Mode (Dashboard)

Run this command to start the system with the real-time web interface. In this mode, you can dynamically adjust the number of sensors, payload sizes, and protocols via the UI.
```bash
make run-dashboard
```
**Access the Dashboard:** Once the containers are healthy, open your browser and navigate to: http://localhost:8080

### 3. Automated Benchmark Mode

To execute the predefined test suite (which automatically iterates through all different configurations) and save raw data for offline analysis:

```bash
make run-bench
```

## Data Analysis

The system generates detailed CSV logs in the results/ folder during benchmark runs. To generate the comparative charts:
 * Ensure you have Python 3 installed.
 * Run the analyzer script:

```bash
python3 analyzer.py
```
This script processes the collected metrics to produce visual comparisons of latency, marshalling time, payload and network overhead.

## Cleanup

To stop the system and remove temporary artifacts, use the following commands:

 * Stop and remove containers:
```bash
make down
```
 * Complete cleanup (containers, logs, and build artifacts):
```bash
make clean
```
