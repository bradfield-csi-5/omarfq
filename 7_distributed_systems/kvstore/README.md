# Distributed Key-Value Store

## Use Case
Serve as a very simple distributed store.
We need:
- Synchronous replication for fault-tolerance and read throughput. If we want to achieve a strong consistency model, replication needs to happen synchronously to avoid having stale data in follower nodes.
    - This can impact write performance since we have to wait
- Avoid replication lag as much as possible
- Idea: what would be the protocol used to distribute read requests  among ndoes? 
    - Randomly distributing them?
    - Using a load balancer (would a need stress test)
    - Using a WAL. Leader would distribute or send the log records across the network to its followers (the drawback with this approach is that is tightly couples the replication strategy with the storage engine since binary formats vary across storage engines).
