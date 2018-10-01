package data

//SimpleMessage is a struct that helps implement a simple
//implementation of the UDP-protocol
type SimpleMessage struct {
	OriginalName  string
	RelayPeerAddr string
	contents      string
}

//GossipPacket is a packet that holds
//onto a SimpleMessage struct
type GossipPacket struct {
	Simple *SimpleMessage
}
