---
title: Architecture
---

[⬅ Back to Home](./) | [Next → Usage](usage.md) | [→ Setup](setup.md)

---

# Usage

Once the Host Monitor service is running:

- **Access the Dashboard:**  
  Open the web dashboard in your browser (e.g., `http://localhost:4200` for dev or your deployed URL).  

- **Real-Time Monitoring:**  
  The UI connects to the backend via **WebSockets**, receiving live updates for each monitored host — no page refresh required.  

- **Historical Data Storage:**  
  Every host check (status, latency, packet loss, timestamp) is **persisted in PostgreSQL** by the backend.  
  This enables future analysis, reporting, and offline queries even if the UI wasn’t connected at the time.  

- **Configure Hosts & Thresholds:**  
  Use the **Slide Panel** to  adjust the latency threshold.  
  Changes are applied instantly without restarting the backend.  

- **Visual Insights:**  
  Review host status cards, uptime pie charts, and latency trend charts in real time.  

- **Export Data:**  
  Use the **Export to CSV** button to download historical uptime and latency data from the live session.  
  For older data, you can query PostgreSQL directly.  

- **Copy Details Quickly:**  
  Click the copy-to-clipboard button on any host card to grab connection details instantly.  

---

[⬅ Back to Setup](setup.md) | [Next → Maintenance](maintenance.md)
