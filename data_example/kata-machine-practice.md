---
title: kata-machine-practice
description: A repository to practice data structures and algoritms following ThePrimeagen's "The last algorithms course you need" from Frontend Masters
label: building
url: https://github.com/zidariu-sabin/kata-machine-practice
updatedat: "2025-10-06T13:38:25Z"
languages:
- name: TypeScript
  color: "#3178c6"
  size: 78.3
- name: JavaScript
  color: "#f1e05a"
  size: 21.7
docfileoid: 0ffebcc
---
## Developed live on twitch
[ThePrimeagen](https://twitch.tv/ThePrimeagen)

### Supported Algorithm
* Insertion sort
* Merge sort
* QuickSort
* Prim's MST (Adjacency List)
* Dijkstra's Shortest Path (Adjacency List)

### Supported Data Structures
* Singly linked list
* Doubly linked list
* Queue
* Stack
* Graph with Adjacency List
* Graph with Adjacency Matrix (untested)

### How It Works

Make sure you have [Node.js](https://nodejs.org/en/) and yarn installed: `npm install --global yarn`

clone the repo and install the dependencies

```bash
yarn install
```

edit the `ligma.config.js` file
```javascript
module.exports = {
    dsa: [
        "InsertionSort",
        "MergeSort",
        "Queue",
        "Stack",
        "QuickSort",
        "DijkstraList",
        "PrimsList",
    ],
}
```

create a day of katas, this will use the list in the `ligma.config.js`.
```bash
yarn generate
```

this will progressively create folders named

```
src/day1
src/day2
...
```

`yarn generate` will also update the `tsconfig.json` and `jest.config` to point
the latest `day` folder via tspaths.  This allows us to avoid updating anything
for testing each day.

#### Testing
```
yarn test
```

I have yet to create a testing strategy for next sets of algorithms, but we
will get there when i cross that bridge.

### Help wanted
A simple way to specify test, thinking something like `tests.json` and `cat
test.json 2> /dev/null` to specify the tests to run.  tests.json wouldn't be
committed.
