package nodes

import (
	"time"

	"github.com/BjornGudmundsson/Peerster/data/hashtable"
)

//timeForLive is how long can pass since a node is declared dead
const timeForLive time.Duration = 10 * time.Second

func (g *Gossiper) CheckLiveTable() {
	for {
		deadNodes := g.LiveTable.DeadNodes(timeForLive)
		if len(deadNodes) != 0 {
			for _, node := range deadNodes {
				g.ChordTable.RemoveNode(node)
			}
			g.RedistributeChunks()
		}
		time.Sleep(5 * time.Second)
	}
}

//RedistributeChunks redistributes the chunks
//to their respective nodes if there has been change in the chord table
func (g *Gossiper) RedistributeChunks() {
	chunks := g.Chunks
	for chunk, d := range chunks {
		indexChunk := hashtable.HashStringInt(chunk)
		pos1, pos2 := g.ChordTable.GetPlaceInChord(indexChunk)
		if pos1 != nil {
			g.GiveChunk(chunk, d, pos1)
		}
		if pos2 != nil {
			g.GiveChunk(chunk, d, pos2)
		}
	}
}
