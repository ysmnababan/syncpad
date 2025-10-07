# Level 2 prioritized checklist (implement in this order)

1. **Core op model + uniqueness**

   * Require every client op to include a stable ID: `opId = clientId + localSeq`.
   * Persist the op **atomically** (write-ahead log or DB transaction) *before* acknowledging the client.
   * Purpose: dedupe retries, reason about idempotence, and guarantee durability.

2. **Server-assigned sequence + optimistic client flow**

   * Server assigns a monotonically increasing **seq** on accept. Client optimistic-applies locally, then reconciles when server ack arrives.
   * Implement client-side in-flight queue: store local ops until acked; when ack arrives, either confirm or transform/rebase the applied op.
   * Purpose: good UX + deterministic global order for replay.

3. **Snapshot + compaction**

   * Periodic snapshot policy (e.g., every N ops or T minutes). Snapshot = materialized document + `lastIncludedSeq`.
   * After snapshot persisted, truncate op log up to `lastIncludedSeq`. Keep snapshot metadata so late joiners request `snapshot + ops > lastIncludedSeq`.
   * Purpose: bounded storage & fast recovery.

4. **Durable storage & recovery**

   * Use a transactional DB or append-only WAL with fsync before ack (sqlite/Postgres/leveldb). On restart: load latest snapshot + replay remaining ops.
   * Implement a clear recovery path and automatic startup checks.

5. **Idempotent apply & client dedupe**

   * Server and client ignore ops whose `opId` already applied. Return idempotent acks.
   * Purpose: robust to retries and network flakiness.

6. **Testing harnesss: concurrency & replay tests**

   * Unit tests to apply randomized interleavings of ops and verify convergence.
   * Replay tests: produce snapshots, then replay only newer ops and assert correctness.
   * Crash-recovery tests: simulate crash after persist/ before persist, restart, and assert no data loss.

7. **Presence & cursors (separate channel)**

   * Keep frequent presence events (cursor) separate from the op log; don’t persist them in the same heavy path. Use ephemeral channels or in-memory broadcast with occasional persistence if needed.

8. **Backups, monitoring, and alerts**

   * Periodic DB backups, snapshot export, and simple logs/metrics (ops/sec, snapshot duration, queue depth).
   * Add health endpoints and a recovery playbook.

9. **Performance & simple scaling**

   * Measure bottlenecks. Add batching (coalesce rapid small ops into larger ops before snapshot or persist) and consider sharding documents across multiple server instances (per-document ownership model).
   * Purpose: scale horizontally by document ownership, without changing core sequencing logic.

---

# Make it upgrade-friendly toward Level 3

* **Keep the sequencing layer modular.** Encapsulate “assign seq & persist” behind a small interface so you can swap it later for a replicated sequencer (Raft) or a CRDT sync provider.
* **Use per-document sequence ranges if you plan to shard.** Per-shard/per-doc sequences avoid global coordination at first.
* **Log format with metadata:** store `opId, clientId, localSeq, serverSeq, timestamp, baseVersion` — this makes it easier to migrate to OT/CRDT mechanisms later.
* **Design clients to tolerate out-of-order acks and rebase operations** — makes future multi-leader or optimistic replication easier.

---

# Low-cost learning wins you should not skip

* Implement the optimistic apply + reconciliation flow — it’s an excellent resume/demo moment.
* Build automated tests simulating concurrent edits — understanding failure modes here is the point of the project.
* Make a short demo recording showing offline edits, reconnect, and clean convergence — recruiters love that.
