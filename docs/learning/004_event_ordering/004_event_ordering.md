# Event Ordering and Causality

Distributed system usually doesn't have a global clock that each system is agree upon. Therefore the order and the causality can't be known exactly.
To tackle that problem, there are some algorithm that can be implemented.

## Physical Clock vc Logical Clock
- Physical clock is a counter for many seconds elapsed
- Logical clock is a counter for numbers of events occured within and across processes.
- Even if physical clock is useful for many things, but it may be inconsistent with causality due to clock skew or network delay.
- That's why we use logical clock where the events is numbered based on the order of when it was occured. Logical clocks instead track `happens-before` relationships (→).

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


## Key takeways
- Logical clock is used for ordering event because the physical clock can be inconsistent across nodes
- Lamport's logical clock can be used to know the `casually total order` but it misses the concurrency info from the processes.
- 