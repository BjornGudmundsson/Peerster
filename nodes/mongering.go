package nodes

import (
	"fmt"
	"math/rand"
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
	usedPeers := make(map[string]bool)
	peers := g.Neighbours
	addr := msg.Addr
	usedPeers[addr] = true
	//Picking a random peer
	peer := peers.GetRandomNeighbour(usedPeers)
	g.Status.ChangeStatus(peer)
	go g.sendRumourMessageToNeighbour(msg.Msg, peer)
	time.Sleep(1 * time.Second)
	var brk bool
	for {
		if brk {
			break
		}
		select {
		case sp := <-g.Status.StatusChannel:
			fmt.Println("Got a message, babe ", sp)
			return
		default:
			coin := rand.Int() % 2
			if coin == 0 {
				fmt.Println("Coin flip said stop")
				g.Status.StopMongering()
				brk = true
			} else {
				fmt.Println("coin flip said continue")
				usedPeers[peer] = true
				nPeer := peers.GetRandomNeighbour(usedPeers)
				g.Status.ChangeStatus(nPeer)
				if nPeer == "" {
					brk = true
					fmt.Println("No more neighbours")
					g.Status.StopMongering()
				}
			}
		}
	}
	fmt.Println("Stopped mongering")
}
