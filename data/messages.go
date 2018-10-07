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
	Status *PeerStatus
}

//MessageHolder Assures a more concurrent access to
//stored messages.
type MessageHolder struct {
	Messages map[string][]RumourMessage
	mux      sync.Mutex
}

//NewMessageHolder returns a pointer to a new MessageHolder
func NewMessageHolder() *MessageHolder {
	return &MessageHolder{
		Messages: make(map[string][]RumourMessage),
	}
}

//AddAMessage is a message bound to a MessageHolder pointer and allows a concurrent
//way to add messages.
func (mh *MessageHolder) AddAMessage(rmsg RumourMessage) {
	mh.mux.Lock()
	n := len(mh.Messages[rmsg.Origin])
	if uint32(n) != rmsg.ID-1 {
		return
	}
	mh.Messages[rmsg.Origin] = append(mh.Messages[rmsg.Origin], rmsg)
	mh.mux.Unlock()
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
