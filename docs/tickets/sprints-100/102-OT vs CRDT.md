# Ticket: Compare OT vs CRDT and Create ADR for Synchronization Strategy

**ID:** TKT-102
**Status:** Done
**Priority:** High
**Planned Score:** 3
**Actual Score:** 3
**Created:** 2025-10-09

---

## 🎯 Goal

Understand the differences between **Operational Transformation (OT)** and **Conflict-free Replicated Data Types (CRDT)**, and decide which model fits our system’s synchronization requirements.
Produce an **ADR document** outlining the decision, rationale, and consequences for distributed data consistency.

---

## ✅ Acceptance Criteria

* [X] Study core concepts and architecture of OT and CRDT
* [X] Summarize comparison: conflict handling, architecture, and consistency guarantees
* [X] Document CRDT family types (state-based, op-based, delta-based) and use cases
* [X] Explain why CRDT is preferred for decentralized or offline-first systems
* [X] Write and store ADR: `ADR-Choosing-CRDT-over-OT.md` under `/docs/decisions`

---

## 📝 Notes

**Resources**

* Shapiro et al., *Conflict-Free Replicated Data Types* (INRIA, 2011)
* Kleppmann, *Designing Data-Intensive Applications*, Chapter 5
* [Automerge Docs](https://automerge.org/)
* [Yjs Docs](https://yjs.dev/)
* [Redis CRDTs](https://redis.io/docs/latest/develop/data-types/crdts/)

**Key Takeaways to Capture**

* CRDTs guarantee **strong eventual consistency** through commutative merge operations.
* OT requires **central coordination and transformation logic**, increasing complexity.
* CRDTs are **more resilient and scalable** in distributed or offline-first architectures.