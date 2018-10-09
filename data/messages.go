package data

import (
	"fmt"
	"sync"
)

//TextMessage just a way to store a
//message in a struct so that it can be
//Serialized and then deserialized.
type TextMessage struct {
	Msg string
}

//SimpleMessage is a struct that helps implement a simple
//implementation of the UDP-protocol
type SimpleMessage struct {
	OriginalName  string
	RelayPeerAddr string
	Contents      string
}

//NewSimpleMessage returns a pointer to a new instance of a
//SimpleMessage structure.
func NewSimpleMessage(ogname, msg, relay string) *SimpleMessage {
	return &SimpleMessage{
		OriginalName:  ogname,
		RelayPeerAddr: relay,
		Contents:      msg,
	}
}

//GossipPacket is a packet that holds
//onto a SimpleMessage struct, a rumour in
//the form of a RumourMessage structu and
//the corresponding status of the node.
type GossipPacket struct {
	Simple *SimpleMessage
	Rumour *RumourMessage
	Status *StatusPacket
}

//MessageHolder Assures a more concurrent access to
//stored messages.
type MessageHolder struct {
	Messages map[string][]RumourMessage
	mux      sync.Mutex
}

//StatusPacket is a holder for PeerStatus messages
type StatusPacket struct {
	Want []PeerStatus
}

//NewMessageHolder returns a pointer to a new MessageHolder
func NewMessageHolder() *MessageHolder {
	return &MessageHolder{
		Messages: make(map[string][]RumourMessage),
	}
}

//AddAMessage is a message bound to a MessageHolder pointer and allows a concurrent
//way to add messages.
func (mh *MessageHolder) AddAMessage(rmsg RumourMessage) bool {
	mh.mux.Lock()
	defer mh.mux.Unlock()
	n := len(mh.Messages[rmsg.Origin])
	if n == 0 {
		if rmsg.ID == 1 {
			mh.Messages[rmsg.Origin] = append(mh.Messages[rmsg.Origin], rmsg)
		}
		return true
	}
	if rmsg.ID <= uint32(n) {
		return false
	}
	if uint32(n) == rmsg.ID+1 {
		mh.Messages[rmsg.Origin] = append(mh.Messages[rmsg.Origin], rmsg)
	}
	return true
}

//GetMessageVector sends a map where the keys are the
//origins of the messages and the values with said key
//are the latest IDs from that origin
func (mh *MessageHolder) GetMessageVector() map[string]int {
	mh.mux.Lock()
	defer mh.mux.Unlock()
	m := make(map[string]int)
	for key, val := range mh.Messages {
		m[key] = len(val)
	}
	return m
}

//NeedMsgs gives a status packet of the messages this node is missing
//according to the given statusPacket
func (mh *MessageHolder) NeedMsgs(sp StatusPacket) StatusPacket {
	want := sp.Want
	messages := mh.GetMessageVector()
	sp2 := StatusPacket{}
	for _, ps := range want {
		msgs := messages[ps.Identifier]
		if ps.NextID > uint32(msgs) {
			p := PeerStatus{
				Identifier: ps.Identifier,
				NextID:     uint32(msgs + 1),
			}
			sp2.Want = append(sp.Want, p)
		}
	}
	return sp2
}

//CreateStatusPacketFromMessageVector takes in a map that represents a
//message vector and returns a StatusPacket made out of said message vector
func CreateStatusPacketFromMessageVector(m map[string]int) *StatusPacket {
	sp := StatusPacket{}
	for key, val := range m {
		ps := PeerStatus{
			Identifier: key,
			NextID:     uint32(val),
		}
		sp.Want = append(sp.Want, ps)
	}
	return &sp
}

//PrintMessagesForOrigin prints all of the messages that have
//come from this origin.
func (mh *MessageHolder) PrintMessagesForOrigin(origin string) {
	mh.mux.Lock()
	messages := mh.Messages[origin]
	fmt.Println("Current messages for origin ", origin)
	for _, msg := range messages {
		fmt.Printf("%v and the ID is %v", msg.Text, msg.ID)
	}
	mh.mux.Unlock()
}
