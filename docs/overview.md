---

## title: Overview

[⬅ Back to Home](./) | [Next → Architecture](architecture.md)

# Overview

**Host Monitor** is a real-time host and service monitoring tool built as part of a skills assessment project.
It is designed to run on Linux, track network host availability, and present live metrics through a modern web dashboard.

## Purpose

The tool continuously pings specified hosts on a given port and interval, collecting metrics that help determine whether a host is **UP** or **DOWN**.

## Key Features

* **Real-time Updates** – Uses WebSockets to push host status and metrics instantly to the UI.
* **Metrics Tracking** – Monitors latency, uptime percentage, and packet loss.
* **Configurable Thresholds** – User-defined latency thresholds to flag degraded performance.
* **Attractive UI** – Angular dashboard with charts, status cards, copy-to-clipboard, and CSV export.
* **Dockerized** – Easy deployment and consistent runtime environment.
* **Extensible Architecture** – Modular Go backend with clear separation of concerns.

## Technology Stack

* **Backend:** Go (Golang)
* **Frontend:** Angular
* **Live Communication:** WebSockets
* **Containerization:** Docker
* **Optional Deployment:** systemd service, Kubernetes manifests

## Intended Audience

This documentation is intended for:

* **Developers** — who want to understand, extend, or maintain the codebase.
* **Operators** — who deploy, configure, and monitor the application in production.

---

[⬅ Back to Home](./) | [Next → Architecture](architecture.md)
