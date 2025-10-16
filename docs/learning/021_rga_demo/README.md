### 🧩 How merges resolve order

When two replicas merge, each operation carries a **unique identifier (replicaID, counter)**.
During merging, replicas compare these IDs and the **causal relationships** recorded in their **version vectors (VV)**.
This allows them to deterministically reconstruct the correct **insertion order** — even if operations were created concurrently — so all replicas eventually converge to the same sequence.

---

### 💀 How tombstones accumulate

In RGA, when an element is deleted, it isn’t physically removed — it’s marked with a **tombstone** flag.
This ensures that future merges don’t accidentally “resurrect” deleted elements that still exist in another replica’s log.
Over time, these tombstones accumulate because each deletion must remain visible to maintain convergence.

---

### ♻️ What I learned about compaction and synchronization

* **Compaction** is needed to reclaim space by safely removing tombstones and redundant history once all replicas have acknowledged those deletions (their VVs confirm everyone has seen them).
* **Synchronization** uses version vectors to exchange only **missing operations**, not the full history — making replication efficient.
  Together, they keep the RGA both **correct (causally consistent)** and **efficient (bounded in size)** over time.
