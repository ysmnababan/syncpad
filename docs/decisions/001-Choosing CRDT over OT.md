# ADR<001>: Choosing CRDT over OT for Distributed Data Synchronization

**Status:** Accepted
**Date:** 2025-10-10
**Context:** Data Synchronization & Conflict Resolution in Distributed Systems

---

## 1. Context

Our system involves **multiple clients or services** that can concurrently modify shared data (e.g., documents, configuration state, metadata).
To maintain **eventual consistency**, we must reconcile concurrent changes without losing updates, even under:

* Network partitions
* Offline edits
* Message delays or reordering

Two candidate approaches for resolving concurrent updates are:

* **Operational Transformation (OT)**
* **Conflict-free Replicated Data Types (CRDT)**

---

## 2. Problem

We need a synchronization mechanism that:

1. Works reliably in **distributed, multi-node environments**.
2. Supports **offline-first editing and merge**.
3. Avoids **central coordination** or single points of failure.
4. Guarantees **eventual consistency** and **no data loss**.
5. Is **easy to reason about** and verifiable.

---

## 3. Option 1: Operational Transformation (OT)

### Core Concept

OT represents edits as operations that must be **transformed** relative to concurrent operations before being applied.
Transformations ensure the same final state across clients by adjusting conflicting positions or intents.

### Characteristics

| Property             | Description                                      |
| -------------------- | ------------------------------------------------ |
| **Consistency**      | Eventual (through transformations)               |
| **Coordination**     | Often centralized (requires total order)         |
| **Complexity**       | High — manual design of transformation functions |
| **Failure Recovery** | Hard — depends on operation history              |
| **Used in**          | Early Google Docs, Etherpad, ShareJS (legacy)    |

**Example:**
Concurrent inserts require index shifting logic to avoid overwriting each other.

```sh
Initial: "Hello"
User A: Insert "X" at position 5
User B: Insert "Y" at position 5

→ Without transformation: order matters, one insert may shift the other
→ With transformation: both inserts preserved
Final: "HelloXY"
```

---

## 4. Option 2: Conflict-free Replicated Data Types (CRDT)

### Core Concept

CRDTs ensure convergence through **mathematically defined merge operations** that are:

* **Commutative** (A + B = B + A)
* **Associative** ((A + B) + C = A + (B + C))
* **Idempotent** (A + A = A)

This allows replicas to update independently and merge at any time without coordination.

### CRDT Families

| Type                        | Description                                | Example Use Case          |
| --------------------------- | ------------------------------------------ | ------------------------- |
| **State-based (CvRDT)**     | Merge full replica state deterministically | Redis counters, Riak sets |
| **Operation-based (CmRDT)** | Broadcast operations that commute          | Yjs, Automerge, Figma     |
| **Delta-based**             | Sync only changed deltas                   | Mobile sync, IoT devices  |

### Common CRDT Data Structures

| CRDT Type                  | Description                 | Use Case                    |
| -------------------------- | --------------------------- | --------------------------- |
| **G-Counter / PN-Counter** | Distributed counters        | Likes, metrics, view counts |
| **OR-Set / 2P-Set**        | Safe add/remove set         | Membership, tags            |
| **LWW-Register**           | Last-write-wins fields      | User preferences, metadata  |
| **RGA / Sequence CRDT**    | Ordered list for text       | Collaborative editing       |
| **Map CRDT**               | Key-value with nested CRDTs | JSON-like data              |
| **Delta CRDT**             | Efficient incremental sync  | Edge / mobile systems       |

**Example**

```sh
Initial counter: 0

Replica A: +1      Replica B: +1

After local updates:
  A = 1
  B = 1

Sync / Merge:
  Merge(A,B) = A + B = 2

Final (both replicas):
  Counter = 2

```

---

## 5. Comparison Summary

| Aspect                | **Operational Transformation (OT)** | **Conflict-free Replicated Data Type (CRDT)** |
| --------------------- | ----------------------------------- | --------------------------------------------- |
| **Approach**          | Transform concurrent ops            | Merge mathematically consistent states        |
| **Consistency**       | Eventual (via transformations)      | Strong eventual consistency                   |
| **Conflict Handling** | Custom per-operation logic          | Automatic convergence                         |
| **Architecture**      | Centralized or ordered              | Fully decentralized                           |
| **Failure Recovery**  | Requires operation history          | Merge is idempotent and safe                  |
| **Offline Support**   | Limited                             | Native                                        |
| **Scalability**       | Lower (central coordination)        | High (no total order)                         |
| **Use Cases**         | Centralized editors                 | Distributed/offline systems                   |

---

## 6. Decision

We **choose CRDT** as the synchronization mechanism for distributed state management.

### Rationale

* **Strong eventual consistency** without central coordination.
* **Fault-tolerant**: replicas can merge after network failures or partitions.
* **Simpler to reason about** formally due to algebraic merge properties.
* **Better scalability** in distributed and peer-to-peer architectures.
* **Offline-first compatibility**, allowing local writes and later synchronization.

OT remains suitable for tightly controlled, centralized, low-latency editors where coordination is guaranteed, but it introduces significant operational and algorithmic complexity in distributed environments.

---

## 7. Consequences

### ✅ Positive

* Simpler distributed design: replicas merge asynchronously.
* No central server dependency.
* High resilience to network issues and client disconnections.
* Natural fit for systems requiring offline editing or edge replication.

### ⚠️ Trade-offs

* Some CRDT types may cause **state growth** (metadata overhead).
* Implementing **complex data types** (e.g., collaborative text) still requires careful design.
* **Causal delivery** or **version vectors** may be needed for efficiency.

---

## 8. References

* Shapiro et al., *Conflict-Free Replicated Data Types*, INRIA (2011)
* Kleppmann, *Designing Data-Intensive Applications*, Ch. 5
* Automerge Docs — [https://automerge.org/](https://automerge.org/)
* Yjs Docs — [https://yjs.dev/](https://yjs.dev/)
* Redis CRDTs — [https://redis.io/docs/latest/develop/data-types/crdts/](https://redis.io/docs/latest/develop/data-types/crdts/)
