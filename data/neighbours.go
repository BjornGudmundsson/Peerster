package data

import (
	"fmt"
	"math/rand"
	"sync"
)

//Neighbours is a struct that allows
//concurrent access to a list of neighbours
type Neighbours struct {
	Neighbours    map[string]bool
	ArrNeighbours []string
	mux           sync.Mutex
}

//AddANeighbour adds a new neighbour to the map
//in a concurrent way.
func (n *Neighbours) AddANeighbour(s string) {
	n.mux.Lock()
	if _, ok := n.Neighbours[s]; !ok {
		n.Neighbours[s] = true
		n.ArrNeighbours = append(n.ArrNeighbours, s)
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
	fmt.Printf("\n")
	n.mux.Unlock()
}

//GetRandomNeighbour is function bound to the neighbours
//struct and gives you a random neighbour that is not
//the same as the given address.
func (n *Neighbours) GetRandomNeighbour(addr map[string]bool) string {
	n.mux.Lock()
	defer n.mux.Unlock()
	l := len(n.ArrNeighbours)
	keys := len(addr)
	if keys >= l {
		return ""
	}
	if l < 2 {
		return ""
	}
	for i := 0; i < l; i++ {
		ran := rand.Int() % l
		if _, ok := addr[n.ArrNeighbours[ran]]; !ok {
			return n.ArrNeighbours[ran]
		}
	}
	return ""
}
