# Distributed Systems

- Aim for parallelism, fault-tolerance, geographical constraint, security/isolation
- Infrastructure that can be seen as distributed system:
  * Storage
  * Communication
  * Computation
- Scalability: ability to increase the infrastructure to improve the performance.
  This can be seen as vertical and horizontal scalability.
- Due to multiple instances running at the same time, distributed systems often introduce
  some faults. So it is very important to build fault-tolerance distributed system.
- Some of core ideas of `fault-tolerance` system:
  * Availability: under some kind of failure, system can still provide services
  * Recoverability: ability to provide service after repair as if nothing happened before.
- Recorevability can be solved by some approaches:
  * Using non volatile memory -> so the data still exist 
  * Replication -> two identical system that can acts as backup when one system is unavailable
- Consistency is a very important quality for the distributed system. Replicated system can have
  the same data for different device but how to make sure the data is consistent across the
  system? Consistency has 2 types:
  * Strong consistency => ensure the operations will synchronize all across the system
    This is very consistent approach, but very expensive to implement.
  * Weak consistency => no guarantee that the system synchronize across the systems right away, 
    but cheaper to implement. This is also known as eventual consistency.


## Replication
Replication can serve several purposes:
 * High availability
   Keeping the system running, even when one machine (or several machines, or an entire datacenter) 
   goes down
 * Disconnected operation
   Allowing an application to continue working when there is a network interruption
 * Latency
   Placing data geographically close to users, so that users can interact with it faster
 * Scalability
   Being able to handle a higher volume of reads than a single machine could handle, 
   by performing reads on replica

Three main approaches to replication:
 * Single-leader replication
   Clients send all writes to a single node (the leader), which sends a stream of data
   change events to the other replicas (followers). Reads can be performed on any
   replica, but reads from followers might be stale. Single-leader replication is popular
   because it is fairly easy to understand and there is no conflict resolution to worry
   about
 * Multi-leader replication
   Clients send each write to one of several leader nodes, any of which can accept
   writes. The leaders send streams of data change events to each other and to any
   follower nodes. This approach can be more robust in the presence of faulty nodes, 
   network interruptions, and latency spikes—at the cost of being harder to reason 
   about and providing only very weak consistency guarantees.
 * Leaderless replication
   Clients send each write to several nodes, and read from several nodes in parallel
   in order to detect and correct nodes with stale data. This approach can be more 
   robust in the presence of faulty nodes, network interruptions, and latency spikes—at
   the cost of being harder to reason about and providing only very weak consistency 
   guarantees.

There are a few consistency models which are helpful for deciding how an application
should behave under replication lag:
  * Read-after-write consistency
    Users should always see data that they submitted themselves.
  * Monotonic reads
    After users have seen the data at one point in time, they shouldn’t later see the
    data from some earlier point in time.
  * Consistent prefix reads
    Users should see the data in a state that makes causal sense: for example, seeing a
    question and its reply in the correct order.

Multi-leader and leaderless replication allow multiple writes to happen con
currently, therefore conflicts may occur. First, we have to determine whether the process
is sequence process or concurrent process. And then we need further study for resolving
conflict this particular case.