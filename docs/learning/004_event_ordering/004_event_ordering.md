# Event Ordering and Causality

Distributed system usually doesn't have a global clock that each system is agree upon. Therefore the order and the causality can't be known exactly.
To tackle that problem, there are some algorithm that can be implemented.

## Physical Clock vc Logical Clock
- Physical clock is a counter for many seconds elapsed
- Logical clock is a counter for numbers of events occured within and across processes.
- Even if physical clock is useful for many things, but it may be inconsistent with causality due to clock skew or network delay.
- That's why we use logical clock where the events is numbered based on the order of when it was occured. Logical clocks instead track `happens-before` relationships (→).

## ⚖️ Causality (Happens-Before Relation)

Defined by Leslie Lamport as **“happens-before” (`→`)**:

`a → b` if:

1. Both in same process, and `a` occurs before `b`, or
2. `a` is sending a message and `b` is receiving it, or
3. There’s a chain: `a → x → ... → b`.

If neither `a → b` nor `b → a`, then `a` and `b` are **concurrent** (`a || b`).

* **Causal relation** = one event could influence another.
* **Concurrent events** = no possible causal link.


## 🔢 Types of Event Ordering

| Type                                | Description                                                        | Detects Concurrency?  | Preserves Causality?   | Example                                           |
| ----------------------------------- | ------------------------------------------------------------------ | --------------------  | --------------------   | ------------------------------------------------- |
| **True Total Order**                | Perfect global timeline — every event has a unique time.           | ❌                    | ✅ (trivially)        | Idealized physical clock (impossible in practice) |
| **Causally Consistent Total Order** | All events ordered, but respects causal order.                     | ❌                    | ✅                    | **Lamport Clock**                                 |
| **Partial Order (Causal Order)**    | Only causally related events are ordered. Concurrent ones are not. | ✅                    | ✅                    | **Vector Clock**                                  |


## Lamport's Clock
- Is a logical clock for imposing `causally consistent total order` of events for multiple processes that .
- `Total order` means all of the events can be viewed as a series of events that happens after each other even if there is concurrent process.
- This clock is useful to know the ordering of the event but can't be used to know whether two process are concurrent or not.
- It guarantees:
  * If `a -> b` then `L(a) < L(b)`
  * But not the reverse. If `L(a) < L(b)`, that doesn't mean that the `a` happens before `b`
- Two concurrent events can have the same Lamport's `timestamp`, because it has to be totally ordered, the might impose a tie-break rule. For example, smaller process ID wins.
- This is the reason why the Lamport's clock is not `true total order`. Even if the `L(a) < L(b)` doesn't necessarly mean that the `a` happens before `b` because those might be concurrent processes but ordered because of the tie-break rule.
- The rule is :
  * If a local event `E` has occured in a process `P`, the counter is increased
  * If a process send a message to another process, sender's counter is increased. Then the increased counter is sent (with the message) to the receiver.
  * If a process receive a message from another process, get the max value between the sender's counter and receiver's last counter. This value is increased by one and add as a new counter for the receiver.
- If you want to the order without missing the concurrency info, please see the [Vector Clock](#vector-clock).
- You can also see the demo for Lamport's Clock [here](./lcdemo/main.go)

```
Process A: a1 --- a2 ---- a3(send 5) --------
                    \
Process B:           ---- b1(receive 5) -----
```

## Vector Clock
- Is a logical clock that tracks the `partial order` of events in a distributed systems.
- Vector clocks let you tell whether two events are:
  * Causally related (a → b)
  * Concurrent (a || b)
- If concurrent is found, later can be processed either by resolving the conflict (e.g. merge, delete, append).
- It guarantees:
  * If `a -> b` then `V(a) < V(b)`
  * If `V(a) < V(b)` then `a -> b`
  * If `V(a) !=< V(b)` and `V(a) !>= V(b)` then, `a || b` (concurrent)
- `V(a) < V(b)` means all of the timestamp of `a` is less than or equal to `b`
  ``` md
    <!-- true -->
    [3 0 0] < [3 1 0]
    [2 0 0] < [3 1 0]
    [5 2 0] < [5 2 2]
  
    <!-- false -->
    [3 1 0] < [5 0 2]
    [5 0 2] < [3 1 0]

  <!-- therefore the [3 1 0] and [5 0 2] is CONCURRENT -->
  ```

- The rule is :
  * If there are 3 process, then each event will has 3 timestamp order
  * If local event is occured: increment the node timestamp while keeping the other timestamp as before
    ``` md
        There are 3 process A,B,C
        The vector becomes => <A,B,C>
      
        <!-- before -->
        In process A = `<0,2,2>`

        <!-- after -->
        In process A = `<1,2,2>`
    ```

  * If a process send a message to another process, sender's vector is increased. Then the increased vector is sent (with the message) to the receiver.
  * If a process receive a message from another process, get the max value between the sender's vector and receiver's last vector. This value is increased by one and add as a new vector for the receiver.
    ``` md
    <!-- before -->
    A: [1 0 0] → [2 0 0] → [3 0 0] → [4 0 2]
    B: [3 1 0]
    C: [2 0 1] → [2 0 2]

    <!-- after A send message to B -->
    A: [1 0 0] → [2 0 0] → [3 0 0] → [4 0 2] → [5 0 2]
    B: [3 1 0]                                         → [5 2 2]
    C: [2 0 1] → [2 0 2]
    ```


## Key takeways
* Logical clock is used for ordering event because the physical clock can be inconsistent across nodes
* **Lamport Clock:**
  "Everyone gets a timestamp that respects causality — but we might invent fake ordering for concurrent events."
* Example : distributed logging system, consensus algo (Paxos, Raft), etc
* **Vector Clock:**
  "We track causality precisely — and can tell when two events happened independently."
* Example : Git, CRDT, Causal Broadcast System (social media feed ordering)