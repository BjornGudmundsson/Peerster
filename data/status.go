package data

//PeerStatus is a struct that that
//keeps the status of a peer.
//The NextID field is the lowest field that
//the sender has not received a message with the
///corresponding ID. The identifier is the name of the node
//where the message originates from.
type PeerStatus struct {
	Identifier string
	NextID     uint32
}
