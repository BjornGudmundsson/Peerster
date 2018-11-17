package nodes

import "fmt"

//HandleSearchRequest handles all incoming search request
//messages coming to the gossiper.
func (g *Gossiper) HandleSearchRequest(msg GossipAddress) {
	fmt.Println("Got a search reply")
}
