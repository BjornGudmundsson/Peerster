package data

import (
	"fmt"
	"sync"
)

//Neighbours is a struct that allows
//concurrent access to a list of neighbours
type Neighbours struct {
	Neighbours map[string]bool
	mux        sync.Mutex
}

//AddANeighbour adds a new neighbour to the map
//in a concurrent way.
func (n *Neighbours) AddANeighbour(s string) {
	n.mux.Lock()
	if _, ok := n.Neighbours[s]; !ok {
		n.Neighbours[s] = true
	}
	n.mux.Unlock()
}

//PrintNeighbours displayes the neighbours of this node
//in a concurrent way such that a neighbour can't be added
//while displaying the current neighbours
func (n *Neighbours) PrintNeighbours() {
	n.mux.Lock()
	fmt.Printf("PEERS ")
	for ip := range n.Neighbours {
		fmt.Printf("%v ", ip)
	}
	n.mux.Unlock()
}
