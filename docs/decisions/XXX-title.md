# ADR <Number>: <Decision Title>

**Date:** YYYY-MM-DD  
**Status:** Accepted | Proposed | Superseded by ADR-XXX  
**Context:** <Why this decision is needed>  
**Decision:** <What choice you made>  
**Consequences:** <Tradeoffs of this choice>  

---

### Example
- Context: We need a conflict resolution algorithm for concurrent edits.  
- Decision: We will use CRDT instead of OT.  
- Consequences: CRDT simplifies merging at scale but has higher storage cost.
