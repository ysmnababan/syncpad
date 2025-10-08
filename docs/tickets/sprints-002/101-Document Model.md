# Ticket: Define and Analyze Initial Collaborative Editor Document Model

**ID:** TKT-101
**Status:** ✅ Done
**Priority:** High
**Planned Score:** 5
**Actual Score:** 5
**Created:** 2025-10-07

---

## 🎯 Goal

Design and deeply understand a baseline architecture for a collaborative editor before choosing the merge algorithm (OT or CRDT).
This includes defining how operations are sent, ordered, stored, compacted, and recovered across clients and replicated servers.

---

## ✅ Acceptance Criteria

* [x] Model designed: includes client/server communication, sequencing, replication, snapshots, and compaction.
* [x] Fault-tolerance and deduplication mechanisms considered (Raft replication, `opId`, and `clientBaseSeq`).
* [x] Understood the purpose of optimistic apply, stale base handling, and compaction.
* [x] Critiqued and documented open problems and undefined parts.
* [x] Identified next focus area: **OT vs CRDT model selection.**

---

## 📝 Notes

### Current Model Summary

* Clients connect via WebSocket to a replicated server cluster (leader via Raft).
* Server assigns a monotonically increasing `serverSeq` (arrival order) to persisted ops.
* Each op carries `opId`, `clientBaseSeq`, and operation payload.
* Server deduplicates ops by `opId` and applies compaction/snapshots periodically.
* If client’s base is stale (`clientBaseSeq < lastIncludedSeq`), server signals `STALE_BASE` for resync.
* Clients perform optimistic local applies and reconcile after receiving server acks.
* Replication ensures leader failover and log durability (details TBD).

### Optional Deeper Dives

* Review how production systems like *Google Docs (OT-based)* and *Figma/Yjs (CRDT-based)* handle synchronization.
* Prepare small visual examples (2 clients editing concurrently) to reason about the merge logic.
* Consider testing frameworks or visualization tools for operation replay later.
