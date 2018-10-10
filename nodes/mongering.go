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

func (g *Gossiper) rumourMongering(msg *GossipAddress) {
	addr := msg.Addr
	gp := msg.Msg
	usedNeighbours := make(map[string]bool)
	usedNeighbours[addr] = true
	randPeer := g.Neighbours.RandomIndexOutOfNeighbours(usedNeighbours)
	fmt.Printf("\n Mongering with %v \n", randPeer)
	usedNeighbours[randPeer] = true
	g.sendMessageToNeighbour(gp, randPeer)
	status := g.Status
	status.ChangeStatus(randPeer)
	var brk bool
	for {
		if randPeer == "" {
			fmt.Println("No more neighbours are left")
			return
		}
		//Sleep while waiting for message
		time.Sleep(1 * time.Second)
		if brk {
			break
		}
		select {
		case msg := <-status.StatusChannel:
			fmt.Println(msg.Msg.Status)
			needMsgs := g.Messages.NeedMsgs(*msg.Msg.Status)
			fmt.Println("These are the messages I need ", needMsgs)
			gp := &data.GossipPacket{
				Status: &needMsgs,
			}
			g.sendMessageToNeighbour(gp, randPeer)
			g.RetrieveMongerMessages()
			randPeer = g.Neighbours.RandomIndexOutOfNeighbours(usedNeighbours)
			usedNeighbours[randPeer] = true
		default:
			coin := rand.Int() % 2
			if coin == 0 {
				fmt.Println("Coin flip said stop")
				g.Status.ChangeStatus("")
				g.Status.StopMongering()
				brk = true
			} else {
				randPeer = g.Neighbours.RandomIndexOutOfNeighbours(usedNeighbours)
				usedNeighbours[randPeer] = true
				fmt.Printf("FLIPPED COIN sending rumor to %v \n", randPeer)
				if randPeer != "" {
					g.sendMessageToNeighbour(gp, randPeer)
				}
			}
		}
	}
}

//RetrieveMongerMessages takes all of the
//coming from the peer this gossiper is mongering with
func (g *Gossiper) RetrieveMongerMessages() {
	mongMsg := g.Mongering
	ch := mongMsg.Ch
	ticker := time.NewTicker(1 * time.Second)
	fmt.Println("Waiting for messages")
	var brk bool
	for {
		if brk {
			return
		}
		select {
		case rm := <-ch:
			fmt.Println("got the message")
			g.Messages.AddAMessage(rm)
			ticker = time.NewTicker(1 * time.Second)
		case <-ticker.C:
			fmt.Println("Ticker expired")
			g.Status.ChangeStatus("")
			g.Status.StopMongering()
			brk = true
		}
	}
}
