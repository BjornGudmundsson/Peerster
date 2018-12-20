package nodes

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/BjornGudmundsson/Peerster/data"
)

//MongererMessages is a struct to help keep track of who the
//application is currently mongering with
type MongererMessages struct {
	Ch chan data.RumourMessage
}

func (g *Gossiper) rumourMongering(rm *data.RumourMessage, addr string) {
	peers := g.Neighbours.GetAllNeighboursWithException(addr)
	replyChannel := make(chan data.GossipPacket)
	fmt.Println("rumour mongering")
	n := len(peers)
	if n == 0 {
		fmt.Println("No peers ")
		return
	}
	t, peer := g.GetPeerNoRumourMonger(peers)
	//I am already rumour mongering with "everyone"
	if !t {
		fmt.Println("Whack")
		return
	}
	g.StatusPeers.AddPeer(peer, replyChannel)
	g.SendRumourMessage(rm, peer)
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case reply := <-replyChannel:
			sp := reply.Status
			uptodate := g.RumourHolder.CheckIfUpToDate(sp)
			if uptodate {
				coinFlip := data.FlipACoin()
				if coinFlip {
					g.StatusPeers.RemovePeer(peer)
					t, peer = g.GetPeerNoRumourMonger(peers)
					if !t {
						return
					}
					g.SendRumourMessage(rm, peer)
				} else {
					fmt.Println("Stopped rumour mongering")
					g.StatusPeers.RemovePeer(peer)
					return
				}
			}
			peerNeeds := g.RumourHolder.GetRumoursPeerNeeds(sp)
			if len(peerNeeds) != 0 {
				//Sending a random rumour that the
				randomRumour := data.GetRandomRumourFromSlice(peerNeeds)
				g.SendRumourMessage(&randomRumour, peer)
			} else {
				//Get the messages that I need if there were no messages that
				//the other peer needs.
				INeed := g.RumourHolder.CheckIfNeedMessages(sp)
				IWant := INeed.Want
				if len(IWant) != 0 {
					g.SendStatusPacket(INeed, peer)
				}
			}
			ticker = time.NewTicker(time.Second)
		case <-ticker.C:
			fmt.Println("Flipping a coin")
			coinFlip := data.FlipACoin()
			if coinFlip {
				g.StatusPeers.RemovePeer(peer)
				t, peer = g.GetPeerNoRumourMonger(peers)
				if !t {
					return
				}
				g.SendRumourMessage(rm, peer)
			} else {
				fmt.Println("Stopped rumour mongering")
				g.StatusPeers.RemovePeer(peer)
				return
			}
		}
	}
}

//GetPeerNoRumourMonger gets a peer that the gossiper is not rumour mongering with.
func (g *Gossiper) GetPeerNoRumourMonger(peers []string) (bool, string) {
	n := len(peers)
	if n == 0 {
		return false, ""
	}
	var gotPeer bool
	r := rand.Int() % n
	peer := peers[r]
	for i := 0; i < 10; i++ {
		hasEntry := g.StatusPeers.HasEntry(peer)
		if hasEntry {
			r = rand.Int() % n
			peer = peers[r]
			continue
		}
		gotPeer = !hasEntry
		break
	}
	return gotPeer, peer
}

//SendRumourMessage sends a rumour message. It is
//an abstraction of sending a rumour message
func (g *Gossiper) SendRumourMessage(rm *data.RumourMessage, addr string) {
	gp := &data.GossipPacket{
		Rumour: rm,
	}
	g.sendMessageToNeighbour(gp, addr)
}

//SendStatusPacket sends a statuspacket to a given neighbour
func (g *Gossiper) SendStatusPacket(sp *data.StatusPacket, addr string) {
	gp := &data.GossipPacket{
		Status: sp,
	}
	g.sendMessageToNeighbour(gp, addr)
}
