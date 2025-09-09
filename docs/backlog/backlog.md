# 📋 Backlog with Acceptance Criteria

### Phase 1 — Core MVP

1. **Design doc v0**
   ✅ A short doc (1–2 pages) with: project goal, core entities (User, Document, Event), initial architecture sketch.

2. **Database setup**
   ✅ Postgres schema created.
   ✅ Tables for users, documents, events, snapshots exist.
   ✅ Can connect and run migrations.

3. **User auth**
   ✅ Signup/login endpoints.
   ✅ Returns JWT with expiry.
   ✅ Middleware rejects requests without valid token.

4. **Document CRUD (basic)**
   ✅ API endpoints to create, read, update metadata.
   ✅ Auth enforced (only logged-in users).

5. **Event append API**
   ✅ Client can POST an edit event.
   ✅ Event is persisted in DB.
   ✅ Returns event ID.

6. **Event replay & snapshot save**
   ✅ API can rebuild document state by replaying events.
   ✅ Snapshot saved automatically after N events.
   ✅ Snapshot used on replay (faster than pure events).

7. **WebSocket server**
   ✅ Clients can connect to a document channel.
   ✅ Sending an edit broadcasts it to all other clients.

8. **Conflict resolution engine**
   ✅ Implement CRDT/OT core functions.
   ✅ Unit tests with concurrent edits converge to same final state.

9. **Basic permissions**
   ✅ Documents have owner/editor/viewer roles.
   ✅ API enforces access control (e.g., viewers cannot edit).

10. **Automated tests (MVP)**
    ✅ Unit tests cover DB, auth, CRDT/OT.
    ✅ Integration test covers: create doc → edit → replay → result correct.

---

### Phase 2 — Essential Production-ish

11. **Introduce message broker**
    ✅ Events are published to Kafka/Redis.
    ✅ Consumer processes and broadcasts events across multiple nodes.

12. **Presence service**
    ✅ Users connecting to a document are stored in Redis.
    ✅ All clients see who is currently online.

13. **Cursor sharing**
    ✅ Each client’s cursor position is broadcast.
    ✅ Other clients see cursor indicators.

14. **Snapshot strategy**
    ✅ Snapshot frequency documented.
    ✅ Load test shows replay time reduced after snapshotting.

15. **Observability: logging**
    ✅ Structured logs (JSON).
    ✅ Request ID included in logs.
    ✅ Logs show request → DB → broker flow.

16. **Observability: metrics**
    ✅ Prometheus endpoint exports counters (requests, errors).
    ✅ Histogram for edit latency.
    ✅ Dashboard shows p95 latency.

17. **Observability: tracing**
    ✅ Jaeger traces exist for: edit event → persist → broadcast.
    ✅ Trace shows spans across DB + broker.

18. **Rate limiting**
    ✅ Redis-based limiter per user.
    ✅ Exceeding limit returns 429 error.

19. **Scaling setup**
    ✅ Run 2+ app nodes behind NGINX.
    ✅ Edits propagate correctly across nodes.

20. **TLS termination**
    ✅ HTTPS endpoint available through NGINX.
    ✅ Self-signed or Let’s Encrypt cert installed.

---

### Phase 3 — Demo polish

21. **Simple frontend demo**
    ✅ Web page with textarea.
    ✅ Multiple users see each other’s edits + cursors live.

22. **Docker Compose**
    ✅ Single `docker-compose up` spins up Postgres, Redis, Kafka, app, and UI.

23. **Demo scenario**
    ✅ Script/test that:

* Logs in two users.
* Both edit same doc.
* Both see edits instantly.

24. **Screencast recording**
    ✅ 2–3 minute video recorded.
    ✅ Shows collaboration + server restart + recovery.

---

### Phase 4 — Portfolio packaging

25. **Architecture diagram**
    ✅ Diagram (PNG/SVG) with boxes: Client, App, DB, Broker, Cache.
    ✅ Arrows show data flow.

26. **Design doc v1**
    ✅ Updated doc (3–5 pages).
    ✅ Covers: CRDT vs OT choice, Kafka vs Redis, snapshot strategy, scaling.

27. **Load test report**
    ✅ Run load test with N concurrent users.
    ✅ Report includes p95/p99 latencies + throughput.
    ✅ Document results in markdown.

28. **README polish**
    ✅ Clear run instructions (`docker-compose up`).
    ✅ Screenshot of demo.
    ✅ Link to screencast + architecture diagram.

29. **Blog-style writeup**
    ✅ Short post (500–1000 words).
    ✅ Explains one hard part (e.g., CRDT).
    ✅ Published on GitHub repo or personal blog.
