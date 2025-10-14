# Ticket: Explore Yjs Behavior and Data Flow

**ID:** TKT-022
**Status:** Todo
**Priority:** High
**Planned Score:** 5
**Actual Score:** <fill in after completion>
**Created:** 2025-10-16

---

## 🎯 Goal

Understand how a **production-grade CRDT library (Yjs)** handles synchronization, garbage collection, and data persistence — to prepare for backend design.
By the end, you should be able to explain:

* How Yjs encodes updates and merges them efficiently
* How it avoids tombstone bloat (compared to RGA)
* What data a backend must persist, relay, or snapshot

---

## ✅ Acceptance Criteria

* [ ] Set up a **minimal Yjs demo** (two browser tabs or Node processes) using `y-websocket` or `y-webrtc`
* [ ] Perform concurrent edits (insert/delete text) and observe convergence
* [ ] Inspect binary updates using:

  * `Y.encodeStateAsUpdate(doc)`
  * `Y.decodeUpdate(update)`
* [ ] Observe how deletions are represented (DeleteSets) and when GC runs (`doc.gc = true/false`)
* [ ] Capture logs or diagrams showing:

  * Update flow (client → server → other clients)
  * What state gets stored/sent
* [ ] Write a reflection comparing Yjs vs RGA:

  * What problems Yjs solved differently
  * How it represents causality and order
  * What backend data model might look like (updates, snapshots, awareness)

---

## 📝 Notes

### Suggested Exploration Steps

1. Use an existing example like [yjs-demos](https://github.com/yjs/yjs-demos) or run a simple `y-websocket` server locally.
2. Create two or three Yjs documents connected to the same room; perform edits concurrently.
3. Log or print Yjs updates in base64 or JSON form.
4. Inspect and document:

   * Structure of updates and DeleteSets
   * When garbage collection removes deleted content
   * Size of updates compared to total document

### Key Intuition Target

> “I understand how Yjs’s CRDT model differs from classic tombstone-based CRDTs, and what data the backend should store to ensure persistence and recovery.”

### Optional Deep Dives

* Try disabling GC (`doc.gc = false`) to see tombstone-like behavior.
* Use `Y.snapshot()` to experiment with compaction.
* Compare binary size of updates vs. a JSON delta.
* Read the Yjs “Internals” documentation for merge and GC algorithms.
