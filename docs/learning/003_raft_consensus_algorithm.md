# Raft Consensus Algorithm

- Consensus is a way to agree on something
- This is very crucial step for achieveing performant and fault tolerant distributed system.
- If a multiple node already aggree on some value, that value will be used as a current value that can be reversed
- There are many consensus algorithm. One of the most common algorithm is Paxos
- This algorithm is fault-tolerant but not very easy to understand and to be implemented
- This is where the Raft consensus comes into play, it has the same functionality but simpler to understand and implement
- Through some mechanism, multiple servers / node can aggree on log of operations, even if some server is crash or unavailable
- Each servers has this replicated logs and each servers's state machine will use this log to perform some operation (e.g. store to db)
- Because each servers has the same orders of replicated logs, the final state of each server will be the same too.

## Raft basic

**Server States**
- Leader: manage all the replicated logs, handle client request
- Follower: receive the replicated logs from leader
- Candidate: attend as a candidate for selecting the next leader

**Terms**
- Is a counter for how many election was held
- New election will increment the term and it is shared across the nodes

**Hearbeat**
- `Leader` sends heartbeat regularly to each `follower` for maintaining a communication and prevent next election
- If any `follower` doesn't receive any `leader` heartbeat within certain timeframe, the next election will be occured.


## Flow
Raft consists of 3 subproblem: **leader election**, **log replication**, **safety**

**Leader Election**
- This process is for determining who is the leader
- All starts as a `follower`
- If a `follower` doesn't receive any communication within certain timeframe, change the state into `candidate`
- The timeframe or timeout is randomized for each `follower` to make sure that there is no split-brain. This happen if many nodes has become the candidates at the exact same time and each of them become a leader. The result is multiple leader and consistency is very difficult to achieve in this case
- After transitioning to a `candidate`, the `terms` will increased by one.
- Each `candidate` will send RequestVote RPC to each node. Also each `candidate` will votes for himself.
- `Candidate` who receive the RequestVote can't vote for every request. There are some condition to be met. For example: the first request vote will be server for each term, if the term of requester is lower than the target node => rejected, if the log of the requester is outdated compared to the target => rejected.
- After a voting, a `candidate` will be a `leader` and the rest will be the `follower`
- The next election period will occur when the leader is faulty or not responding for certain time.


**Log Replication**
- After the `leader` is elected, it starts handling client requests.
- When the leader receives a client request, it first **logs** (appends) the command to its own log, with the current term.
- The leader sends AppendEntries RPCs to each follower, including fields: previous log index & term, the new entries, and its current commitIndex.
- Followers perform consistency check (prevLogIndex/prevLogTerm). If that fails, they reject; if succeeds, they append new entries (dropping conflicting entries).
- Leader waits for replies; once a **majority** of followers have stored the entry, the leader marks the entry as *committed*.
- The leader applies committed entries in order to its own state machine.
- In subsequent AppendEntries RPCs (or even heartbeats), the leader includes the commitIndex so followers learn of commit.
- Followers, upon seeing commitIndex greater than theirs, apply the entries up to commitIndex to their state machines, in order.
- Then the leader responds to the client (once the entry is committed and applied locally).
- If a follower’s log is behind or has conflicts, the leader uses backtracking: uses nextIndex / matchIndex, retries with adjusted prevLogIndex/Term until follower accepts.

**Safety**
- The primary objective of safety is to guarantee that all the non-failed server/node executed the same log entries at the same order
- But there might be the case where the server is unavailable during a commit performed by a leader. If this server is elected, the replicated log might be missing some commit and therefore introduce inconsistency into the system
- To prevent this, there is mechanism called election restriction.
- Each server will reject (not giving a vote) to a `candidate` which has outdated log compared to the target server.
- Under this mechanism, a voter (server receiving a RequestVote from a candidate) will only grant its vote if:
  1. The candidate’s term in the RequestVote is ≥ the voter's currentTerm.
  2. The candidate has not already voted for someone else in this same term.
  3. The candidate’s log is at least as up-to-date as the voter's log. Concretely, this is determined by comparing:
     * The candidate’s last log term vs the voter’s last log term; if the candidate’s is greater, it’s more up to date.
     * If the last terms are equal, then compare the last log indices: the candidate must have at least as high an index.
- This mechanism ensure that all the elected leader has the latest log, therefore ensure the safety of the replication.

## How Raft avoid split-brain?

**Randomized election timeouts**:  Each follower uses a timeout that’s randomly chosen from a range (e.g. 150-300 ms). This makes it unlikely multiple followers time out (become candidates) at *exactly* the same time. If one times out earlier, it starts election, wins (if it gets majority), becomes leader, sends heartbeats, which resets election timers on others. This breaks symmetry. 

**Majority / quorum voting** : In election, a candidate must receive votes from a majority of the nodes. If the network is partitioned, only a partition that contains a majority can elect a leader. Smaller partitions cannot, so they don’t get a leader and thus don’t risk split brain. 

**Term numbers and term updates** : Raft maintains a `term` integer that increases on each new election. Messages (RPCs) carry term. If a candidate or leader sees a higher term in a message, it steps down (becomes follower). This ensures that stale leaders are stopped once they learn of a more recent leader election. 

**At most one vote per term per server** : Each server votes at most once per term. So even if multiple candidates request votes, the voter grants only one. This limits the possibility that two candidates both think they have majority. Combined with majority requirement, ensures only one can win. 

**Log up-to-date check in RequestVote**  : When a candidate requests votes, the voter checks whether the candidate’s log is *at least as up-to-date* as its own (via the checks of last log term / last log index). This prevents a candidate with stale log (missing committed entries) from being elected leader. This helps avoid conflicting state diverging. While not directly “split brain” prevention (i.e. two leaders), it's safety for log consistency, which also helps the system behave correctly. 

**Heartbeat mechanism** : Once a leader is elected, it sends periodic heartbeats (AppendEntries RPCs with no new entries) to followers. These reset their election timeouts. So followers don’t start elections while the leader is operating. This prevents multiple simultaneous leader elections when a valid leader is still functioning.