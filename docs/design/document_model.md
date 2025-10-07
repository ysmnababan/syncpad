# One-sentence summary

You’re using a centralized sequencer model: clients connect over WebSocket, send local ops (with client-side counters/ids), the server assigns a global monotonically increasing sequence number (server-seq) for each accepted op, persists ops, issues snapshots + compaction periodically, and clients optimistic-apply then reconcile on ack.

# What your model gets right (strengths)

* **Deterministic global order:** server Seq gives a single canonical ordering so replicas can converge if ops are applied deterministically.
* **Durability path available:** persisting ops before ack + snapshots/compaction is compatible with Level-2 durability.
* **Good UX option:** optimistic local apply + ack reconciliation supports responsive editing.
* **Low initial operational complexity:** single sequencer is simple to implement and great for a learning/demo/resume project.

# Remaining correctness issues (why it can still diverge / be surprising)

* **Numeric positions remain fragile.** Clients send `pos` computed against local state; server ordering can shift positions. Without transform/rebase semantics or a position model that survives reordering, inserts/deletes land wrong.
* **Ambiguity about what a server sequence orders.** You must be explicit: serverSeq orders *committed application order*. Any other interpretation invites mismatches.
* **Tie-breaking / intention preservation missing.** Total order resolves deterministically but may not match user intent (two concurrent inserts at same logical spot still need deterministic tie-break or metadata to preserve intent).
* **Client optimistic vs server reality mismatch.** If clients apply optimistically, they need a clear, implementable way to adjust when serverAck implies transforms — otherwise UI jumps or lost-ops happen.

# Fault-tolerance & operational pitfalls

* **Single sequencer = single point of failure / throughput bottleneck.** Acceptable for Level-2 but you’ll need persistence + restart/recovery to mitigate downtime/data loss.
* **Crash / persistence race:** issuing sequence numbers without durable persistence before acknowledging can lead to gaps, duplicates, or reissued seqs after crashes.
* **Pre-allocation / batch-seq semantics are dangerous if unclear.** Handing ranges of seqs to clients or pre-allocating without robust recovery semantics creates unused holes or conflicting assumptions.
* **Sharding/scale complexity:** global seqs don’t scale across shards; moving to per-doc or per-shard sequences changes ordering semantics and requires careful design if cross-doc consistency is needed.

# Storage / lifecycle issues

* **Op log growth:** without snapshots + truncation you’ll hit storage and replay costs.
* **Snapshot/compaction correctness:** you must record `lastIncludedSeq` and ensure late-joining clients only need snapshot + ops after that seq. Otherwise replays can resurrect or reapply stale operations.
* **Idempotency & dedupe:** need stable `opId` (clientId + localSeq) to dedupe retries and prevent duplicate application.

# Security / practical concerns

* **Authenticate & authorize ops:** ensure ops carry user identity and are validated before persisting/ordering.
* **Prevent malicious replays:** combine opIds + serverSeq + auth checks to avoid attackers reintroducing or faking ops.

# Short, non-prescriptive hints (do these next; no full algorithms)

* Use stable `opId = clientId + clientSeq` for idempotence and dedupe.
* Make serverSeq assignment atomic and durable before acknowledging clients.
* Record explicit metadata with each op: `opId, clientId, clientBaseVersion, serverSeq, timestamp`.
* Implement snapshots with `lastIncludedSeq` and truncate ops ≤ that seq only after snapshot persistence.
* Keep presence/cursor messages out of the durable op path (ephemeral transport).
* Implement optimistic apply with an in-flight queue; when serverAck arrives, rebase/transform the in-flight ops against intervening serverSeq ops.

# Concrete short checklist you can copy into your notes (next three hands-on steps)

1. Add `opId` and persist op atomically on server before ack.
2. Add snapshot + metadata `lastIncludedSeq` and truncate up to that seq.
3. Implement client in-flight queue: optimistic apply → store local op with `opId` → send → on ack match by `opId` and reconcile.

# Final take / how to use this

This model is a solid Level-2 starting point: it teaches you sequence assignment, persistence, snapshots, optimistic UI, and common failure modes (all resume-worthy). But keep in mind the two big knobs you’ll face later: (A) how to represent positions so they survive reordering (transform/rebase or CRDT positional ids), and (B) how far you’ll take fault tolerance (single durable sequencer vs distributed consensus). Tackle (A) next when you’re ready to address correctness under concurrency.