package data

import (
	"fmt"
	"sync"
)

//RumoursFromPeer checks if
//it has seen this rumour before
type RumoursFromPeer struct {
	Count    uint32
	Messages []RumourMessage
}

//NewRumoursFromPeers creates a new RumoursFromPeer struct.
func NewRumoursFromPeers() *RumoursFromPeer {
	return &RumoursFromPeer{
		Count:    0,
		Messages: make([]RumourMessage, 0),
	}
}

//IsNew checks if this is a new rumour from this peer
func (rfp *RumoursFromPeer) IsNew(rm RumourMessage) bool {
	id := rm.ID
	if id <= rfp.Count {
		return false
	}
	return true
}

//AddRumour adds a rumour to the struct if its ID is one
//greater than the current counter. It does not add empty
//string to the struct. It will return true if the rumour
//was added
func (rfp *RumoursFromPeer) AddRumour(rm RumourMessage) bool {
	id := rm.ID
	if id == rfp.Count+1 {
		rfp.Messages = append(rfp.Messages, rm)
		rfp.Count = rfp.Count + 1
		return true
	}
	return false

}

//GetRumoursFromIndex returns a list of all the rumours from a given index - 1
func (rfp *RumoursFromPeer) GetRumoursFromIndex(nxt uint32) *RumourMessage {
	if nxt > rfp.Count {
		//Returning an empty rumour message if it is higher
		return nil
	}
	nxtIndex := nxt - 1
	return &rfp.Messages[nxtIndex]
}

//RumourHolder holds all of the
//rumourmessages this peer has received.
type RumourHolder struct {
	mux             sync.Mutex
	Rumours         map[string]*RumoursFromPeer
	MessagesInOrder []RumourMessage
}

//NewRumourHolder returns a new
//empty RumourHolder struct
func NewRumourHolder() *RumourHolder {
	return &RumourHolder{
		Rumours:         make(map[string]*RumoursFromPeer),
		MessagesInOrder: make([]RumourMessage, 0),
	}
}

//IsNew checks if the rumour is new to the rumour holder
func (rh *RumourHolder) IsNew(rm RumourMessage) bool {
	rh.mux.Lock()
	defer rh.mux.Unlock()
	src := rm.Origin
	if _, ok := rh.Rumours[src]; ok {
		rfp := rh.Rumours[src]
		return rfp.IsNew(rm)
	}
	return true
}

//AddRumour adds a rumour to the RumourHolder struct
func (rh *RumourHolder) AddRumour(rm RumourMessage) {
	rh.mux.Lock()
	defer rh.mux.Unlock()
	src := rm.Origin
	if rfp, ok := rh.Rumours[src]; ok {
		added := rfp.AddRumour(rm)
		if added {
			rh.MessagesInOrder = append(rh.MessagesInOrder, rm)
		}
	} else {
		rfp := NewRumoursFromPeers()
		added := rfp.AddRumour(rm)
		fmt.Println(*rfp)
		rh.Rumours[src] = rfp
		if added {
			rh.MessagesInOrder = append(rh.MessagesInOrder, rm)
		}
	}
}

//GetMessagesInOrder returns the messages that have been received
//in the order that they were received.
func (rh *RumourHolder) GetMessagesInOrder() []RumourMessage {
	return rh.MessagesInOrder
}

//CreateStatusPacket creates a status packet from the
//rumours in the rumourholder struct
func (rh *RumourHolder) CreateStatusPacket() *StatusPacket {
	var want []PeerStatus
	for name, rfp := range rh.Rumours {
		nxt := rfp.Count + 1
		ps := PeerStatus{
			Identifier: name,
			NextID:     nxt,
		}
		want = append(want, ps)
	}
	return &StatusPacket{
		Want: want,
	}
}

//GetRumoursPeerNeeds gets all the rumours that this peer has that the peer
//that send the statuspacket does not have.
func (rh *RumourHolder) GetRumoursPeerNeeds(sp *StatusPacket) []RumourMessage {
	want := sp.Want
	tempRumours := make([]RumourMessage, 0)
	for _, peer := range want {
		src := peer.Identifier
		nxt := peer.NextID
		rfp, ok := rh.Rumours[src]
		if !ok {
			continue
		}
		rumour := rfp.GetRumoursFromIndex(nxt)
		if rumour == nil {
			continue
		}
		tempRumours = append(tempRumours, *rumour)
	}
	return tempRumours
}

//CheckIfNeedMessages checks the messages that this peer needs that are present
//in the peerster.
func (rh *RumourHolder) CheckIfNeedMessages(sp *StatusPacket) *StatusPacket {
	var want []PeerStatus
	for _, ps := range sp.Want {
		src := ps.Identifier
		nxt := ps.NextID
		rfp, ok := rh.Rumours[src]
		if !ok {
			p := PeerStatus{
				Identifier: src,
				NextID:     1,
			}
			want = append(want, p)
			continue
		}
		count := rfp.Count
		if count < nxt {
			p := PeerStatus{
				Identifier: src,
				NextID:     count + 1,
			}
			want = append(want, p)
			continue
		}
	}
	return &StatusPacket{
		Want: want,
	}
}

//CheckIfUpToDate checks if the rumour holder is up to date with
//this statuspacket
func (rh *RumourHolder) CheckIfUpToDate(sp *StatusPacket) bool {
	want := sp.Want
	fmt.Println(want)
	counterMap := make(map[string]uint32)
	for _, ps := range want {
		counterMap[ps.Identifier] = ps.NextID
	}
	n1 := len(counterMap)
	n2 := len(rh.Rumours)
	if n1 != n2 {
		return false
	}

	for name, rfp := range rh.Rumours {
		nxt, ok := counterMap[name]
		if !ok {
			return false
		}
		if nxt != rfp.Count+1 {
			return false
		}
	}
	return true
}
