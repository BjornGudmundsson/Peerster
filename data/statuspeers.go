package data

import "sync"

//StatusPeers is a struct to keep state of
//which peer the gossiper is gossiping with.
//It can allow for  concurrent access and is used to
//lead a status packet to the appropriate process
type StatusPeers struct {
	mux   sync.Mutex
	Peers map[string]chan GossipPacket
}

//NewStatusPeers returns a new empty
//StatusPeers struct.
func NewStatusPeers() *StatusPeers {
	return &StatusPeers{
		Peers: make(map[string]chan GossipPacket),
	}
}

//PassPacketToProcess pass the packet from a peer to the appropriate process that
//is waiting for it.
func (sp *StatusPeers) PassPacketToProcess(name string, packet GossipPacket) {
	sp.mux.Lock()
	defer sp.mux.Unlock()
	if _, ok := sp.Peers[name]; ok {
		sp.Peers[name] <- packet
	}
}

//HasEntry checks if the name of a given peer is in the struct
func (sp *StatusPeers) HasEntry(name string) bool {
	sp.mux.Lock()
	defer sp.mux.Unlock()
	_, ok := sp.Peers[name]
	return ok
}

//AddPeer adds a peer to the statuspeers struct
func (sp *StatusPeers) AddPeer(name string, ch chan GossipPacket) {
	sp.mux.Lock()
	defer sp.mux.Unlock()
	if _, ok := sp.Peers[name]; !ok {
		sp.Peers[name] = ch
	}
}

//RemovePeer removes the entry of a peer in
//the statuspeers struct.
func (sp *StatusPeers) RemovePeer(name string) {
	sp.mux.Lock()
	defer sp.mux.Unlock()
	delete(sp.Peers, name)
}
