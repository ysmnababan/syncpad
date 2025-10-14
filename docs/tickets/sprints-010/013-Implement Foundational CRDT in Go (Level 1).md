# Ticket: Implement Foundational CRDTs in Go (Level 1)

**ID:** TKT-013
**Status:** Done
**Priority:** High
**Planned Score:** 3
**Actual Score:** 6
**Created:** 2025-10-10

---

## 🎯 Goal

Build intuition on how **commutative replication** and **eventual consistency** work in CRDTs by implementing and simulating simple state-based types (G-Counter, PN-Counter, and OR-Set) in Go.

By the end of this ticket, you should *feel confident* explaining:

* How CRDTs achieve convergence without central coordination
* Why operations can arrive in any order and still yield the same state
* How to merge distributed state safely

---

## ✅ Acceptance Criteria

* [X] Implement at least **two Level-1 CRDTs** in Go: G-Counter, PN-Counter, or OR-Set
* [X] Simulate **at least 2–3 replicas** (using goroutines or functions) exchanging state and merging periodically
* [X] Observe and record behavior when updates arrive in different orders
* [X] Write a short summary explaining what “eventual convergence” *feels like* in practice
* [X] Capture 2–3 takeaways relevant to backend design (e.g., merge semantics, idempotency, metadata growth)

---

## 📝 Notes

**Suggested steps:**

* Start with in-memory structs and local merging logic
* No networking required — simulate with local goroutines or simple function calls
* Optional: Add randomized message delay to visualize eventual consistency
* Don’t worry about deletes or compaction yet
* Resources to skim:

  * *CRDTs Illustrated* ([https://crdt.tech](https://crdt.tech))
  * Martin Kleppmann’s blog “A Critique of Operational Transformation” (background motivation)

**Intuition target:**

> “I understand how state-based CRDTs reach the same final value no matter what order operations arrive, and I can visualize merge behavior in Go code.”

