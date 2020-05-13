# Peer-to-peer Privacy-Preserving Asynchronous Range Queries
CSCI 339 Final Project

## Introduction
We describe a peer-to-peer privacy-preserving location-based exposure notification (read: digital contact tracing) system.
To preserve privacy, we hash location data and use a system of rotating node-IDs.
The system notifies users every 6-12 hours of any potential exposure events.

### Challenges
The primary challenges of building this system are as follows:
 - using location data in a privacy-preserving manner, and
 - continuously monitoring distance between moving nodes. 

## Overview
Our system consists of two phases: (1) planting and (2) tracking.  The planting phase "plants" tracks that will later be tracked in the tracking phase.  The system runs continuously, taking snapshots of user position every so often (we use the constant `GAP` to denote the time between snapshots). 
We define two other constants:
let `THRESH` be the minimum exposure time that tracks an exposure notification, 
and let `K` be the distance threshold for infection transmission.

We partition the network domain into a matrix of disjoint cells, where each cell has length and width `K`. We can then refer to each cell `c` via a pair of coordinates `name(c) = (x,y)` denoting the latitude and longitude of the top left corner of the cell.

We use a distributed hashtable (DHT) to store hashed location information.  Each entry in the table consists of a key-value pair `<k,v>` where `k = hash("(x,y)")` for some x,y, and `v` is a list of `(node-id, interval(a,b))` pairs.

### Planting
The planting phase is performed in batches every `GAP/2` time steps.
Whenever a node `n` moves to a new cell `c` at time `t`, we store the cell and time locally via `log`.
Later, in batches, the nodes "plant" trackers by updating the DHT via `plant`.

```
func log(node, c, t):
  let cells = map(closed-neighborhood(node), (c => name(c))

  for cell in cells:
    add interval(t) to node.log[cell], merging to existing intervals if possible
```

```
func plant(node):
  for cell, intervals in node.log:
    if cell not in DHT:
      DHT[cell] = empty
    add (id(node.seed), interval) pairs to DHT[cell]
```

(Use caching to store the node IP that stores the `cell` key to avoid having to repeatedly search for the same cell.)
(Use an interval tree for `DHT[cell]`?)

We assume the existence of a function `id(seed)`, which outputs a new pseudorandom number generated using `seed` every time
it is called.

### Tracking
The tracking phase is performed in batches every `GAP/2` time steps for each node using its stored data.

```
func track(node):
  for cell, interval_1 in node.log:
    let seen = all (id, interval_2) pairs for all intervals in DHT[cell] that overlap interval_1
    for (id, interval_2) in seen:
      node.seen[id] += len(intersect(interval_1, interval_2))
```

(TODO: How can we make the interval checking faster?  Sorting intervals by starting point? Using an interval tree?)

When an infection is diagnosed, the IDs emitted by the patient are broadcast through the network.
All nodes receive the broadcast and check whether they have seen the IDs, and if they have, for how long.

```
func receive-broadcast(node, broadcast-list):
  let sum = 0
  for id in broadcast-list:
    sum += node.seen.get(id, default=0)
    if sum > THRESH:
      alert user of exposure
```

Nodes remove all local records after 30 days.

## References

[1] P. Gladames, Ying Cai, Kihwan Kim. 
    A Generic Platform for Efficient Processing of Spatial Monitoring Queries in Mobile Peer-to-Peer Networks

[2] DP3T

[3] Stoica et al. 
    Chord

[4] Tanin et al. 
    An Efficient Nearest Neighbor Algorithm for P2P Settings
    
