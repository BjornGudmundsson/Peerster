package data

import (
	"fmt"
	"sync"
)

//RoutingTable is a struct that has a map
//of all the nodes and the "next" hop to reach
//that node. It updates itself in a concurrent way
type RoutingTable struct {
	Table map[string]string
	mux   sync.Mutex
}

//UpdateRoutingTable is a function that updates the routing table
//in a concurrent way
func (rt *RoutingTable) UpdateRoutingTable(origin string, hop string) {
	rt.mux.Lock()
	rt.Table[origin] = hop
	fmt.Printf("\n DSDV %v %v \n", origin, hop)
	rt.mux.Unlock()
}
