package data

import (
	"fmt"
	"sync"

	"github.com/BjornGudmundsson/Peerster/data/transactions"
)

//TextMessage just a way to store a
//message in a struct so that it can be
//Serialized and then deserialized.
type TextMessage struct {
	Dst      string
	Msg      string
	File     string
	Request  string
	Keywords string
	Budget   uint64
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
	Simple         *SimpleMessage
	Rumour         *RumourMessage
	Status         *StatusPacket
	PrivateMessage *PrivateMessage
	DataRequest    *DataRequest
	DataReply      *DataReply
	SearchReply    *SearchReply
	SearchRequest  *SearchRequest
	TxPublish      *transactions.TxPublish
	BlockPublish   *transactions.BlockPublish
}

//MessageHolder Assures a more concurrent access to
//stored messages.
type MessageHolder struct {
	Messages     map[string][]RumourMessage
	messageArray []RumourMessage
	mux          sync.Mutex
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

//CheckIfMsgIsNew takes in a rumour message and checks if the
//messageholder has a copy of that message already.
func (mh *MessageHolder) CheckIfMsgIsNew(rmsg RumourMessage) bool {
	mh.mux.Lock()
	defer mh.mux.Unlock()
	og := rmsg.Origin
	id := rmsg.ID
	originMsgs := mh.Messages[og]
	n := len(originMsgs)
	if n == 0 {
		return true
	}
	//The messages should be in order.
	//We are working with that assumption.
	latestMsg := originMsgs[n-1]
	if id > latestMsg.ID {
		return true
	}
	return false
}

//AddAMessage is a message bound to a MessageHolder pointer and allows a concurrent
//way to add messages.
func (mh *MessageHolder) AddAMessage(rmsg RumourMessage) {
	mh.mux.Lock()
	og := rmsg.Origin
	id := rmsg.ID
	msgs := mh.Messages[og]
	n := len(msgs)
	if id == uint32(n+1) {
		mh.messageArray = append(mh.messageArray, rmsg)
		msgs = append(msgs, rmsg)
	}
	mh.Messages[og] = msgs
	mh.mux.Unlock()
}

//GetMessageVector sends a map where the keys are the
//origins of the messages and the values with said key
//are the latest IDs from that origin
func (mh *MessageHolder) GetMessageVector() map[string]int {
	mh.mux.Lock()
	defer mh.mux.Unlock()
	messages := mh.Messages
	m := make(map[string]int)
	for key, arr := range messages {
		//The array is has up to this message
		//I am assuming that it is in order.
		//Fx that ID nr i comes before Id nr i + 1
		m[key] = len(arr) + 1
	}
	return m
}

//NeedMsgs gives a status packet of the messages this node is missing
//according to the given statusPacket
func (mh *MessageHolder) NeedMsgs(sp StatusPacket) StatusPacket {
	mh.mux.Lock()
	defer mh.mux.Unlock()
	sp2 := StatusPacket{}
	want := sp.Want
	messages := mh.Messages
	for _, ps := range want {
		id := ps.Identifier
		nxt := ps.NextID
		msgs, ok := messages[id]
		if !ok {
			ps := PeerStatus{
				Identifier: id,
				NextID:     uint32(1),
			}
			sp2.Want = append(sp2.Want, ps)
			continue
		}
		n := len(msgs)
		if uint32(n) >= nxt {
			continue
		}
		ps := PeerStatus{
			Identifier: id,
			NextID:     uint32(n + 1),
		}
		sp2.Want = append(sp2.Want, ps)
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

//PrintMessages prints all of the messages
//in a specific format.
func (mh *MessageHolder) PrintMessages() {
	for _, val := range mh.Messages {
		for _, msg := range val {
			fmt.Printf(" %v ", msg)
		}
		fmt.Println()
	}
}

//GetStatusPacketFromVector takes in a map that represent a
//message vector and turns it into a StatusPacket
func GetStatusPacketFromVector(m map[string]int) StatusPacket {
	sp := StatusPacket{}
	for key, val := range m {
		ps := PeerStatus{
			Identifier: key,
			NextID:     uint32(val),
		}
		sp.Want = append(sp.Want, ps)
	}
	return sp
}

func (mh *MessageHolder) GetMessageString() string {
	var s string
	for _, val := range mh.Messages {
		for _, rm := range val {
			temp := fmt.Sprintf("Origin %v ID %v Content %v", rm.Origin, rm.ID, rm.Text)
			s += temp
		}
		s += "\n"
	}
	return s
}

func (mh *MessageHolder) GetMessagesInOrder() []RumourMessage {
	return mh.messageArray
}
