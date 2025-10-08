## Ticket: Study Event Ordering & Causality

**ID:** TKT-004
**Status:** Done
**Priority:** Medium
**Planned Score:** 5
**Actual Score:** 5
**Created:** 2025-10-03

### 🎯 Goal

Understand how distributed systems reason about order and causality of events without global clocks.

### ✅ Acceptance Criteria

* [ ] (Optional) Read Lamport’s paper: *Time, Clocks, and Ordering of Events* ([PDF](https://pdos.csail.mit.edu/6.824/papers/times.pdf))
* [X] Learn vector clocks from at least one tutorial or lecture
* [X] Implement or run a tiny simulation (Lamport clock or vector clock) to see event ordering
* [X] Write a short note on the difference between Lamport clocks and vector clocks
* [X] Capture 2–3 takeaways (why causality matters for OT/CRDT)

### 📝 Notes
Tutorial videos:
  - [Lamport's Clock](https://www.youtube.com/watch?v=mo8OPP5FCTg)
  - [Logical Time](https://www.youtube.com/watch?v=x-D8iFU1d-o)
Event ordering is critical for reasoning about concurrent updates.