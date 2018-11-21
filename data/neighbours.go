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

//GetAllNeighboursWithException is a function bound to a Neighbours struct
//that returns the list of all neighbours with the exception of a specified neighbour
func (n *Neighbours) GetAllNeighboursWithException(addr string) []string {
	temp := make([]string, 0)
	for key := range n.Neighbours {
		if key != addr {
			temp = append(temp, key)
		}
	}
	return temp
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

//RandomIndexOutOfNeighbours returns a random neighbour that has not been used
func (n *Neighbours) RandomIndexOutOfNeighbours(used map[string]bool) string {
	l := len(n.ArrNeighbours)
	if len(used) > l {
		return ""
	}
	for {
		ran := rand.Int() % l
		neighbour := n.ArrNeighbours[ran]
		if _, ok := used[neighbour]; !ok {
			return neighbour
		}
	}
}
