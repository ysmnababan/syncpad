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
