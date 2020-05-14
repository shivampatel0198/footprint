# Peer-to-peer Privacy-Preserving Asynchronous Range Queries
CSCI 339 Final Project

## Introduction
Many of the more prominent digital contact-tracing protocols (like DP3T) use local communication between devices via Bluetooth.  We seek to build a protocol that uses location information instead in a privacy-preserving manner.
We do this by hashing location data and using a system of rotating time-sensitive node IDs.

## Overview
We partition the network domain into a matrix of disjoint square cells, where each cell has side length `DIST_THRESH`. We refer to each cell `C` by the pair of coordinates `name(C) = (x,y)` where `x` is the latitude and `y` the longitude of `C`'s top left corner.

At fixed intervals (every `STEP` minutes), each node takes a location snapshot recording which cell it is in (`log()`). This generates a chronological list of (cell, interval) pairs that describes the node's position over time.  *This information is only stored locally*.

We use a distributed hashtable (DHT) to securely store global location information.
The DHT maps from cells to lists of (ID, interval) pairs. 
This represents a heatmap of which cells have been occupied when, and by whom.
Because the location values are hashed, and each node continually generates temporary IDs, the DHT cannot be used to determine users' locations directly; nor can it be used to determine a particular node's path through space. 
Nodes update the DHT with their locally logged location information in batches, via `push()`.
When a patient tests positive, their logged location information is used to identify temporary IDs of other nodes occupying cells near the patient's node.  These IDs are then broadcast through the network: listening nodes determine whether these IDs match any IDs that they previously emitted. This is used to measure a node's exposure.
If exposure exceeds `EXPOS_THRESH`, then a notification is triggered.

## Core Functionality
Our system consists of a few core functions: `log()`, `push()`, `broadcast()`, and `receive()`.

### `log()`
Nodes continually log location data and generated temporary IDs locally.
The internal log is structured as a dictionary mapping from cells to lists of intervals. Used temporary IDs are stored in `node.ids`.

(Note: We store the node's most recent location to skip logging when the node is not moving.)

```
func log(node, cell, t):
  if cell == node.prevCell:
    node.hasNotMoved = true
    return

  // Track the cells around cell that are close enough to the node
  // to constitute transmission events
  let cells = get-neighbors(cell) 

  for cell in cells:
    add interval(t) to node.log[cell], merging into existing intervals if possible,
    taking node.hasNotMoved into account
```

### `push()`
Nodes periodically push their location logs to the DHT.
Nodes retain these logs locally, as they are needed in the case of infection.

```
func push(node):
  for (cell, interval) in node.log:
    if cell not in DHT:
      initialize DHT[cell]

    let id = id(node.seed)
    add (id, interval) to DHT[cell]
    add id to node.ids // Keep track of all emitted IDs
```

(TODO: Use caching to store the node IP that stores the `cell` key to avoid having to repeatedly search for the same cell.  Consider using an interval tree for DHT[cell].)

The function `id(node)` generates a time-sensitive pseudorandom number using the node's private seed:
```
func id(node):
  return hash(node.seed + currentTime)
```

### `broadcast()`
When a patient tests positive, we run `broadcast()` to notify nodes that have had contact events with the patient's node.

```
func broadcast(infected-node):

  // Collect the nodes that have had contact events with the infected node
  contact-events = []
  for (cell, interval) in infected-node.log:
    for (id, interval') in DHT[cell] where interval' overlaps interval:
      add (id, overlap(interval, interval')) to contact-events

  // Propagate affected node IDs through network
  propagate(neighbors(infected-node), contact-events, propagate)
```
(Note: Be careful to avoid flooding the network with repeated messages.)

### `receive()`
When nodes receive a `(id, overlap)` broadcast, they respond by checking whether they emitted an ID matching `id` and adjusting their exposure count accordingly.  If total exposure exceeds `EXPOS_THRESH`, then the node notifies the user.

```
func receive(node, id, overlap):

  if id is in node.ids:
    node.exposure += overlap

  if exposure > EXPOS_THRESH:
    alert user of exposure
```

Nodes remove all local records after 2 weeks.

## Comparison to other techniques

1. Decentralized Privacy-Preserving Proximity Tracing (DP3T)

2. Spatial Monitoring Queries (Galdames)

## System Properties

1. Although the protocol relies on location information, *all location information is 
   stored locally*. Any location information that is exposed to other nodes is hashed 
   first.
   
2. Third parties cannot use global data to track an individual user.

## References

[1] P. Galdames, Ying Cai, Kihwan Kim. 
    A Generic Platform for Efficient Processing of Spatial Monitoring Queries in Mobile Peer-to-Peer Networks

[2] DP3T

[3] BlueTrace

[4] Stoica et al. 
    Chord

[5] Tanin et al. 
    An Efficient Nearest Neighbor Algorithm for P2P Settings
    
