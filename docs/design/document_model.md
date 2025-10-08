# Latest document model (concise)

**Topology & transport**

* Clients connect to a cluster of servers over WebSocket.
* Servers run a replicated sequencer (leader elected via Raft). Leader assigns canonical `serverSeq` to persisted ops.

**Client lifecycle**

* On connect: client fetches latest snapshot metadata (`lastIncludedSeq`) + recent ops (`> lastIncludedSeq`) and sets `currentServerSeq`.
* Client does **optimistic apply** locally when the user edits and enqueues the operation in an *in-flight queue*.

**Operation message**

* Client → Server per op:

  * `opId` (clientId + clientLocalSeq),
  * `clientId`,
  * `clientBaseSeq` (serverSeq the client last knew),
  * `op` (insert/delete payload, incl. positions or payload format).
* Server persists op atomically (WAL/DB) and assigns `serverSeq` (arrival order at leader), then broadcasts ack/broadcast including `serverSeq` and either `appliedOp` or `interveningOps`.

**Dedup & retries**

* Server uses `opId` to detect duplicates; if a duplicate is received it returns the earlier `serverSeq`/ack (idempotency).

**Snapshots & compaction**

* Server periodically creates snapshots (materialized doc + `lastIncludedSeq`) and truncates/compacts op log up to `lastIncludedSeq`.
* If a client sends an op with `clientBaseSeq < lastIncludedSeq`, server returns a `STALE_BASE` response instructing client to re-sync snapshot + ops.

**Replication & availability**

* Sequencer is replicated (Raft style); leader assigns `serverSeq`. Followers replicate persisted ops. Leader election handles failover (details TBD).

**Misc**

* Presence/cursor updates are kept ephemeral (separate from durable op pipeline).
* Clients reconcile optimistic in-flight ops when server acks arrive (rebase/transform logic run client or handled by server).

---

# Concise list of unanswered/unspecified concerns (things to address later)

**Correctness & semantics**

* **Position model undefined:** raw numeric `pos` is still fragile under concurrent edits — who transforms indices and what are the exact transform rules? (OT vs CRDT choice not resolved.)
* **Definition of `serverSeq` semantics:** it’s arrival order at leader, but is that *ordering of persisted ops after transform* or *ordering of original arrivals*? Need exact canonical semantics to avoid ambiguity.

**Transform & reconciliation**

* **Who performs transforms?** Server-side transforms reduce client work but centralize complexity; client-side rebase requires clients to run transform logic. Not specified.
* **In-flight multi-op causality:** how multiple local in-flight ops are causally related and how to rebase them in order.

**Replication & leader semantics**

* **Atomicity & durability at leader:** must persist before ack; not specified how to ensure no gaps on leader crash.
* **Leader failover effects:** what happens to seq continuity and uncommitted ops during leader change? (raft commit/replication quorum behavior needs exact handling.)
* **Pre-allocation / holes:** if you ever pre-allocate seq ranges to speed throughput, how do you handle unused holes on failover?

**Snapshots / compaction under replication**

* **Coordination of snapshots across replicas:** when leader compacts/truncates, followers must agree on `lastIncludedSeq`; need protocol for snapshot installation and consistent truncation.
* **Client resync semantics:** exact flow for `STALE_BASE` (should client re-send ops or re-create them after resync?) not specified.

**Storage & recovery**

* **Write-ahead logging & fsync policy:** how do you guarantee persisted ops will survive crash before acknowledging clients? Performance tradeoffs not decided.
* **Crash-after-assign race:** if a seq is assigned but not durably replicated before ack, crash can cause contradictions — not yet defined.

**Scalability & sharding**

* **Per-doc vs global sequencing:** global seqs don’t scale; plan for per-document/per-shard seqs or ownership model not defined.
* **Cross-shard edits or multi-doc transactions:** not addressed.

**Compaction GC & tombstones**

* **Deletes & garbage collection:** how are deletions/Ghosts/tombstones represented and compacted (especially for CRDT-like approaches)? Not defined.
* **Late-joining replicas & tombstone retention:** retention policy and GC invariants missing.

**Undo / history / provenance**

* **Undo semantics in multi-user context** (intention-preserving undo) not specified; transforms or special metadata required.

**Performance & UX**

* **InterveningOps size bounds:** how to prevent huge digests for lagging clients (throttle, snapshots, delta compression) not specified.
* **Backpressure and flow control:** how server protects from high fanout or slow clients is unspecified.

**Security & integrity**

* **Auth / authorization / op signing:** who is allowed to submit ops and how to prevent malicious replays/injection not defined.
* **Auditability:** how transformed ops preserve provenance for auditing.

**Testing & upgrades**

* **Test harnesses:** specific test cases for leader crash, duplicate opId, stale-base, reorderings not yet enumerated.
* **Versioning & migration:** op & snapshot versioning strategy for future changes not defined.

---

# Short, prioritized things to decide next (so you can move forward)

1. Choose OT vs CRDT for the position/merge model (this affects almost everything else).
2. Specify exact serverSeq semantics (arrival vs transformed-application) and persist-before-ack policy.
3. Define stale-base resync flow concretely (what client does on `STALE_BASE`).
4. Decide per-document sequencing or global sequencing (scale implications).