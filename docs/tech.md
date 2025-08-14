---
title: Technology Choices
---


[⬅ Back to Home](./) | [Start → Overview](overview.md) | [Architecture](architecture.md) | [Usage](usage.md) | [Maintenance](maintenance.md) | [Future Improvements](future_improvements.md)


# Technology Choices

## Why Go (Golang) 

* **Performance and Efficiency** – Go is compiled, statically typed, and known for fast startup times and low memory usage, making it well-suited for a network monitoring tool.
* **Concurrency Model** – Go’s goroutines and channels provide lightweight, built-in concurrency, ideal for handling multiple host pings in parallel.
* **Company Alignment** – Go (Golang) is already a core part of the company’s technology ecosystem, ensuring seamless integration with existing systems and developer familiarity.
* **Personal Growth** – I wanted to learn a modern, low-level language, and Go provided a good blend of performance and simplicity.

## Why Angular for Frontend

* **Rich Component Ecosystem** – Angular provides a structured framework with built-in tooling, which accelerates building a responsive dashboard.
* **Strong TypeScript Support** – Helps maintain type safety and catch errors early.
* **Reactive UI Updates** – Works seamlessly with WebSockets to update charts and metrics in real-time.

## Why WebSockets over REST

* **Real-Time Communication** – REST would require polling, which is inefficient for continuous monitoring. WebSockets push data instantly when status changes.
* **Lower Latency** – Persistent connection reduces the overhead of repeated HTTP requests.
* **Better UX** – Users see host status updates immediately without refreshing or waiting for intervals.

## Why Docker

* **Portability** – Ensures the application runs consistently across environments.
* **Simplified Deployment** – Easy to spin up both backend and frontend in isolated containers.
* **Integration Ready** – Works with CI/CD pipelines, Kubernetes, and cloud platforms.

 [Start → Architecture](architecture.md) | [Setup](setup.md) | [Usage](usage.md) | [Maintenance](maintenance.md) | [Future Improvements](future_improvements.md)

---

<sub>© 2025 Host Monitor • <a href="https://github.com/BenjaminBatte/host-monitor">GitHub Repo</a></sub>
