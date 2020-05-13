# Peer-to-peer Privacy-Preserving Asynchronous Range Queries
CSCI 339 Final Project

## Introduction
We describe a peer-to-peer privacy-preserving location-based exposure notification (read: digital contact tracing) system.
To preserve privacy, we hash location data and use a system of rotating node-IDs.
The system notifies users every 6-12 hours of any potential exposure events.

### Challenges
The primary challenges of building this system are using location data in a privacy-preserving manner, and monitoring the distance between moving nodes without local communication. 

## Overview
Our system consists of two phases: (1) planting and (2) triggering.  
The planting phase updates a global hash table that will later be used to track exposure events in the triggering phase.  
The system runs continuously, taking snapshots of user position every so often (we use the constant `STEP` to denote the time between snapshots). 
We define two other constants:
`EXPOS_THRESH` is the minimum exposure time that triggers an exposure notification, 
and `K` is the distance threshold for infection transmission.

We partition the network domain into a matrix of disjoint square cells, where each cell has side length `K`. 
We refer to each cell `c` via a pair of coordinates `name(c) = (x,y)` where `x` is the latitude and `y` the longitude of `c`'s top left corner.

We use a distributed hashtable (DHT) to store hashed location information.  
Each entry in the table consists of a key-value pair `<k,v>` 
where `k = hash("(x,y)")` for a consistent hashing function `hash` and some reals x,y, 
and `v` is a list of `(id, interval(a,b))` pairs.

### Planting
The planting phase is performed in batches every `STEP/2` time steps.
Whenever a node `n` moves to a new cell `c` at time `t`, we record the event by locally logging the time and all of the 8-neighbors of cell `c` (i.e., any cell that shares a corner vertex with `c`) via `log`.

```
func log(node, c, t):
  let cells = map(closed-neighborhood(node), c => name(c))

  for cell in cells:
    add interval(t) to node.log[cell], merging to existing intervals if possible
```

Later, in batches, nodes "plant" triggers by updating the DHT via `plant`.

```
func plant(node):
  for cell, intervals in node.log:
    if cell not in DHT:
      initialize DHT[cell]
    add (id(node.seed), interval) pairs to DHT[cell]
```

(TODO: Use caching to store the node IP that stores the `cell` key to avoid having to repeatedly search for the same cell.  Use an interval tree for DHT[cell]?)

We assume the existence of a function `id(seed)`, which generates a new pseudorandom number using `seed` every time
it is called.

### Triggering
The triggering phase is performed in batches every `STEP/2` time steps for each node using its stored data.

```
func trigger(node):
  for cell, interval_1 in node.log:
    let seen = all (id, interval_2) pairs for all intervals in DHT[cell] that overlap interval_1
    for (id, interval_2) in seen:
      node.seen[id] += len(intersect(interval_1, interval_2))
```
(How can we make the interval checking faster?  Using an interval tree?)

When an infection is diagnosed, the IDs emitted by the patient are broadcast through the network.
All nodes receive the broadcast and check whether they have seen the IDs, and if they have, for how long.
If they have seen them for more than `EXPOS_THRESH` time steps, then we alert the user.

```
func receive-broadcast(node, broadcast-list):
  let sum = 0
  for id in broadcast-list:
    sum += node.seen.get(id, default=0)
    if sum > EXPOS_THRESH:
      alert user of exposure
```

Nodes remove all local records after 30 days.

## Comparison to other techniques

1. Decentralized Privacy-Preserving Proximity Tracing (DP3T)

2. Spatial Monitoring Queries (Galdames)

## References

[1] P. Galdames, Ying Cai, Kihwan Kim. 
    A Generic Platform for Efficient Processing of Spatial Monitoring Queries in Mobile Peer-to-Peer Networks

[2] DP3T

[3] Stoica et al. 
    Chord

[4] Tanin et al. 
    An Efficient Nearest Neighbor Algorithm for P2P Settings
    
