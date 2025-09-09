# 📝 syncpad

**syncpad** is a collaborative real-time editor backend.  
It’s inspired by tools like Google Docs, and demonstrates how to build a backend system that supports **live editing**, **conflict resolution**, and **event replay**.

---

## 🚀 Features (MVP goals)
- User authentication (JWT-based)
- Document CRUD operations
- Event sourcing (append-only log of edits)
- Snapshot + replay for fast document recovery
- Real-time collaboration via WebSockets
- Conflict resolution with CRDT/OT
- Basic permissions (owner, editor, viewer)

---

## 📂 Project structure
```

syncpad/
backend/       # Backend source code
frontend/      # Demo UI
infra/         # Infra configs (docker-compose, etc.)
docs/          # Backlog, tickets, ADRs, spikes, designs

````

---

## 🛠️ Getting started
### Prerequisites
- Docker + Docker Compose  
- Node.js (if running frontend locally)  

### Run everything
```bash
docker-compose up
````

### Access

* Backend API → [http://localhost:4000](http://localhost:4000)
* Demo UI → [http://localhost:3000](http://localhost:3000)

---

## 📖 Documentation

All docs are under [`docs/`](./docs):

* [Backlog](./docs/backlog)
* [Tickets](./docs/tickets)
* [Decisions (ADR)](./docs/decisions)
* [Spikes](./docs/spikes)
* [Design docs](./docs/design)

---

## 🎯 Why this project?

This project was created as a **learning exercise** and **portfolio project** to demonstrate:

* Building a distributed backend with **event sourcing** and **real-time communication**
* Applying **CRDT/OT algorithms** for conflict resolution
* Documenting design decisions and tradeoffs (ADR, spikes)
* Setting up production-like tooling (logging, metrics, snapshots, scaling)

---

## 📜 License

MIT
