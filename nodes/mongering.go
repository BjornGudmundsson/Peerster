package nodes

import (
	"fmt"
	"sync"
	"time"
)

//Mongerers is a struct to help keep track of who the
//application is currently mongering with
type Mongerers struct {
	Peers map[string]bool
	mux   sync.Mutex
}

//AddMongerer allows for a concurrent way to
//add a mongerer to the Mongerers struct
func (m *Mongerers) AddMongerer(addr string) {
	m.mux.Lock()
	m.Peers[addr] = true
	m.mux.Unlock()
}

//DelMongerer removes a mongerer from the
//Mongerers struct in a concurrent way.
func (m *Mongerers) DelMongerer(addr string) {
	m.mux.Lock()
	m.Peers[addr] = false
	m.mux.Unlock()
}

//HasMongerer allows for a concurrent way to check
//if a Mongerers struct has a specific address
func (m *Mongerers) HasMongerer(addr string) bool {
	m.mux.Lock()
	defer m.mux.Unlock()
	return m.Peers[addr]
}

func (g *Gossiper) rumourMongering(msg *GossipAddress) {
	peers := g.Neighbours
	addr := msg.Addr
	//Picking a random peer
	peer := peers.GetRandomNeighbour(addr)
	fmt.Println(peer)
}

func (g *Gossiper) sendRumourToPeer(msg *GossipAddress, peer string) {
	rumour := msg.Msg
	g.sendRumourMessageToNeighbour(rumour, peer)
	status := g.Status
	status.ChangeStatus(peer)
	time.Sleep(1 * time.Second)
	select {
	case sp := <-status.StatusChannel:
		fmt.Println(sp.Want)
	default:
		fmt.Println("Timeout")
	}
}
