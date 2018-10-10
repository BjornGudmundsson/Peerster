package nodes

import (
	"sync"

	"github.com/BjornGudmundsson/Peerster/data"
)

//Status is a struct that keeps track of which peer
//the gossiper is waiting for a message from
type Status struct {
	IP            string
	IsMongering   bool
	StatusChannel chan GossipAddress
	mux           sync.Mutex
}

//GossipAddress is just a way for me to send both
//the address and GossipPacket through the same channel.
type GossipAddress struct {
	Addr string
	Msg  *data.GossipPacket
}

//ChangeStatus allows for a concurrent way
//to make sure that not many activites are trying to change
//it at the same time
func (s *Status) ChangeStatus(peer string) {
	s.mux.Lock()
	s.IP = peer
	s.IsMongering = true
	s.mux.Unlock()
}

//StopMongering sets the IsMongering
//to false meaning that the mongering
//has stopped.
func (s *Status) StopMongering() {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.IsMongering = false
}

//GetIP returns the IP of the struct
//concurrent way
func (s *Status) GetIP() string {
	s.mux.Lock()
	defer s.mux.Unlock()
	return s.IP
}
