package nodes

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/BjornGudmundsson/Peerster/data"
	"github.com/BjornGudmundsson/Peerster/data/hashtable"
)

//HandleChunkStoreRequest handles chunk store request structs
func (g *Gossiper) HandleChunkStoreRequest(msg GossipAddress) {
	csr := msg.Msg.ChunkStoreRequest
	dst := csr.Destination
	if dst != g.Name {
		g.ForwardChunkStoreRequest(csr)
		return
	}
	hash := csr.Hash
	data := csr.Data
	hexData := string(data)
	//This is a very liberal system. we just accept any chunk
	g.Chunks[hash] = hexData
	g.SendChunkStoreReply(hash, csr.Src, g.Name, true)
}

//SendChunkStoreRequest send a chunk store request to target destination
func (g *Gossiper) SendChunkStoreRequest(hash, dst, src string, d []byte) {
	nxtHop, ok := g.RoutingTable.Table[dst]
	if !ok {
		//Not in the routing table
		fmt.Println("Not in the routing table")
		return
	}
	csr := &hashtable.ChunkStoreRequest{
		Hash:        hash,
		Destination: dst,
		Src:         src,
		Data:        d,
	}
	gp := &data.GossipPacket{
		ChunkStoreRequest: csr,
	}
	g.sendMessageToNeighbour(gp, nxtHop)
}

//ForwardChunkStoreRequest handles forwarding ChunkStoreRequest
func (g *Gossiper) ForwardChunkStoreRequest(msg *hashtable.ChunkStoreRequest) {
	dst := msg.Destination
	nxtHop, ok := g.RoutingTable.Table[dst]
	if !ok {
		//Somethign went wrong. Not in the routing tables
		return
	}
	gp := &data.GossipPacket{
		ChunkStoreRequest: msg,
	}
	g.sendMessageToNeighbour(gp, nxtHop)
}

//SendChunkStoreReply sends a chunk store reply to the destination
func (g *Gossiper) SendChunkStoreReply(hash, dst, src string, r bool) {
	nxtHop, ok := g.RoutingTable.Table[dst]
	if !ok {
		//Something went wrong here. Not in the routing table.
		return
	}
	reply := &hashtable.ChunkStoreReply{
		Hash:        hash,
		Destination: dst,
		Reply:       r,
		Src:         src,
	}
	gp := &data.GossipPacket{
		ChunkStoreReply: reply,
	}
	g.sendMessageToNeighbour(gp, nxtHop)
}

//SpreadMetaFile spreads the metafile to their respective nodes
func (g *Gossiper) SpreadMetaFile(data []byte) {
	n := len(data)
	hash := sha256.New()
	mfh := hash.Sum(data)
	hexHash := hex.EncodeToString(mfh)
	p := hashtable.HashStringInt(hexHash)
	p1, p2 := g.ChordTable.GetPlaceInChord(p)
	if p1 != nil {
		g.GiveChunk(hexHash, hex.EncodeToString(data), p1)
	}
	if p2 != nil {
		g.GiveChunk(hexHash, hex.EncodeToString(data), p2)
	}
	for i := 0; i < n; i = i + 32 {
		j := i + 32
		chunk := data[i:j]
		hexChunk := hex.EncodeToString(chunk)
		val := g.Chunks[hexChunk]
		b := hashtable.HashStringInt(hexChunk)
		pos1, pos2 := g.ChordTable.GetPlaceInChord(b)
		if pos1 != nil {
			g.GiveChunk(hexChunk, val, pos1)
		}
		if pos2 != nil {
			g.GiveChunk(hexChunk, val, pos2)
		}
	}
}

//GiveChunk sends a chunk to the appropriate node
func (g *Gossiper) GiveChunk(chunk, d string, b *big.Int) {
	dstNode := g.ChordTable.GetNodeAtPosition(b)
	if dstNode == g.Name {
		g.Chunks[chunk] = d
		return
	}
	identifier := dstNode + "-" + chunk
	ch := make(chan *hashtable.ChunkStoreReply)
	g.ReplyHandler.AddProcess(identifier, ch)
	counter := 0
	//We will wait up to 10 seconds for a reply to arrive
	ticker := time.NewTicker(2 * time.Second)
	Data := []byte(d)
	g.SendChunkStoreRequest(chunk, dstNode, g.Name, Data)
	for {
		fmt.Println("Bjorninnnn")
		select {
		case reply := <-ch:
			fmt.Println("got reply", reply.Reply, reply.Destination)
			g.ReplyHandler.RemoveProcess(identifier)
			return
		case <-ticker.C:
			fmt.Println("ticker expired")
			if counter == 5 {
				g.ReplyHandler.RemoveProcess(identifier)
				return
			}
			ticker = time.NewTicker(2 * time.Second)
			counter = counter + 1
			g.SendChunkStoreRequest(chunk, dstNode, g.Name, Data)
			time.Sleep(1 * time.Second)
		}
	}
}

//HandleChunkStoreReply handles ChunkStoreReply
func (g *Gossiper) HandleChunkStoreReply(msg *data.GossipPacket) {
	reply := msg.ChunkStoreReply
	dst := reply.Destination
	hash := reply.Hash
	if dst != g.Name {
		g.SendChunkStoreReply(hash, dst, reply.Src, true)
		return
	}
	g.ReplyHandler.PassToProcess(reply)
}
