# Ticket: Prototype RGA (Replicated Growable Array) in Go

**ID:** TKT-021
**Status:** Todo
**Priority:** High
**Planned Score:** 5
**Actual Score:** <fill in after completion>
**Created:** 2025-10-14

---

## 🎯 Goal

Implement and experiment with an **RGA (Replicated Growable Array)** CRDT in Go to understand how sequence CRDTs achieve convergence using tombstones and causal ordering.
By the end, you should intuitively understand:

* How inserts and deletes propagate without coordination
* Why tombstones exist and how they affect metadata growth
* How concurrent inserts are resolved into a consistent global order

---

## ✅ Acceptance Criteria

* [ ] Implement a minimal **RGA** supporting `Insert(afterID, value)` and `Delete(id)`
* [ ] Represent elements with `(id, value, tombstone, prevID)` and merge logic between replicas
* [ ] Simulate **at least 2–3 replicas** performing concurrent operations (insert/delete) and merging
* [ ] Verify **eventual convergence**: all replicas reach identical visible sequences
* [ ] Log/visualize tombstones and element order during merges
* [ ] Write a short summary explaining:

  * How merges resolve order
  * How tombstones accumulate
  * What you learned about compaction and synchronization

---

## 📝 Notes

### Suggested Implementation Steps

1. Define element struct and replica state (map + ordered list).
2. Implement insert/delete operations with causal references (`prevID`).
3. Write merge function:
   * Union all elements
   * Rebuild order using causal links
   * Ignore tombstones for visible text
4. Simulate replicas exchanging states out of order (randomized merge order).
5. Add logs to visualize sequence and tombstone count.

### Suggested Experiments

* **Concurrent insert:** Two replicas insert at the same position. Observe final order.
* **Concurrent delete + insert:** One deletes, one inserts after the same ID. Observe resolution.
* **Tombstone growth:** Run repeated insert/delete cycles and record how tombstones accumulate.

### Intuition Target

> “I can explain how a simple sequence CRDT converges using causal links and tombstones — and why production systems (like Yjs) need GC and compact delta formats.”

### Optional Deep Dives

* Add version vectors for merge optimization.
* Compare RGA merge vs. naive list merge.
* Write reflection: “What part of RGA would be hardest to scale?”
