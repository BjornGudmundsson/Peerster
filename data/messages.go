package data

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
