package nodes

import (
	"fmt"
	"log"
	"net"

	"github.com/BjornGudmundsson/Peerster/data"
	"github.com/BjornGudmundsson/Peerster/data/blockchain"
	"github.com/BjornGudmundsson/Peerster/data/transactions"
	"github.com/dedis/protobuf"
)

//Gossiper is an instance of a node running the peerster protocol.
//It has the UDP address of it self and its active UDP connection.
type Gossiper struct {
	address               *net.UDPAddr
	conn                  *net.UDPConn
	Name                  string
	Neighbours            *data.Neighbours
	Messages              *data.MessageHolder
	Counter               *data.Counter
	Status                *Status
	Mongering             MongererMessages
	enPeer                *EntropyPeer
	UIPort                int
	RoutingTable          *data.RoutingTable
	PrivateMessageStorage *data.PrivateMessageStorage
	Files                 map[string]data.MetaData
	Chunks                map[string]string
	MetaFileHashes        data.MetaFileHashes
	StateFileFinder       data.StateFileFinder
	HandlerDataReplies    data.HandlerDataReplies
	RecentRequest         data.RecentRequests
	ChunkToPeer           *data.ChunkToPeer
	FoundFileRepository   data.FoundFileRepository
	BlockChain            *blockchain.BlockChain
	TransactionBuffer     *transactions.TransactionBuffer
	RumourHolder          *data.RumourHolder
	StatusPeers           *data.StatusPeers
}

//NewGossiper is a function that returns a pointer
//to a new instance of a Gossiper with an active
//UDP connection to the specified address.
func NewGossiper(address, name string, neighbours []string, p int) *Gossiper {
	udpaddr, e := net.ResolveUDPAddr("udp4", address)
	if e != nil {
		log.Fatal(e)
	}
	udpconn, e := net.ListenUDP("udp4", udpaddr)
	if e != nil {
		log.Fatal(e)
	}
	m := data.SliceToBoolMap(neighbours)
	conN := &data.Neighbours{
		Neighbours:    m,
		ArrNeighbours: neighbours,
	}
	status := &Status{
		IP:            "",
		IsMongering:   false,
		StatusChannel: make(chan GossipAddress),
	}
	mong := MongererMessages{
		Ch: make(chan data.RumourMessage),
	}
	counter := &data.Counter{}
	routingTable := &data.RoutingTable{
		Table: make(map[string]string),
	}
	tempMap := make(map[string][]string)
	privStorage := data.PrivateMessageStorage(tempMap)
	files := make(map[string]data.MetaData)
	chunks := make(map[string]string)
	mfh := data.NewMetaFileHashes()
	sff := data.NewStateFileFinder()
	hdr := data.NewHandlerDataReplies()
	recentrequests := data.NewRecentRequests()
	ctp := data.NewChunkToPeer()
	ffr := data.NewFoundFileRepository()
	bc := blockchain.NewBlockChain()
	txb := transactions.NewBuffer()
	rh := data.NewRumourHolder()
	sp := data.NewStatusPeers()
	return &Gossiper{
		Name:                  name,
		address:               udpaddr,
		conn:                  udpconn,
		Neighbours:            conN,
		Messages:              data.NewMessageHolder(),
		Counter:               counter,
		Status:                status,
		Mongering:             mong,
		enPeer:                &EntropyPeer{},
		UIPort:                p,
		RoutingTable:          routingTable,
		PrivateMessageStorage: &privStorage,
		Files:                 files,
		Chunks:                chunks,
		MetaFileHashes:        mfh,
		StateFileFinder:       sff,
		HandlerDataReplies:    hdr,
		RecentRequest:         recentrequests,
		ChunkToPeer:           ctp,
		FoundFileRepository:   ffr,
		BlockChain:            bc,
		TransactionBuffer:     txb,
		RumourHolder:          rh,
		StatusPeers:           sp,
	}
}

func (g *Gossiper) sendMessageToNeighbour(msg *data.GossipPacket, addr string) {
	udpaddr, e := net.ResolveUDPAddr("udp4", addr)
	if e != nil {
		log.Fatal(e)
	}
	packetbyte, e := protobuf.Encode(msg)
	if e != nil {
		log.Fatal(e)
	}
	_, e = g.conn.WriteToUDP(packetbyte, udpaddr)
	if e != nil {
		log.Fatal(e)
	}
}

//ReceiveMessages listens for incoming messages coming
//in on the UDP connection for this node.
func (g *Gossiper) ReceiveMessages() {
	conn := g.conn
	gossipChannel := make(chan GossipAddress)
	go g.delegateMessages(gossipChannel)
	for {
		buffer := make([]byte, 8*1024)
		gp := &data.GossipPacket{}
		_, addr, _ := conn.ReadFromUDP(buffer[:])
		protobuf.Decode(buffer, gp)
		go g.Neighbours.AddANeighbour(addr.String())
		gAddress := GossipAddress{
			Addr: addr.String(),
			Msg:  gp,
		}
		gossipChannel <- gAddress
	}
}

//This is the threshold for
//for the number of hits
//required to be satisfied
//with the  search
const threshold uint64 = 2

//This is the maximum of requests
//I'll send if the budget was not specified
const maxBudget int = 32

//ClientMessageReceived is a function bound to a pointer to the Gossiper struct.
//It enables a listening on the port specified in the function parameters
//For packages coming from the client.
func (g *Gossiper) ClientMessageReceived(port int) {
	addr := g.address.IP
	fullAddr := fmt.Sprintf("%v:%v", addr.String(), port)
	udpAddr, _ := net.ResolveUDPAddr("udp4", fullAddr)
	conn, _ := net.ListenUDP("udp4", udpAddr)
	packet := make([]byte, 1024)
	temp := &data.TextMessage{}
	for {
		n, _, e := conn.ReadFromUDP(packet)
		if e != nil {
			fmt.Println(n, e.Error())
		}
		protobuf.Decode(packet, temp)
		if temp.File != "" {
			g.FileHandling(*temp)
			continue
		}
		if temp.Keywords != "" {
			go g.HandleClientSearchRequests(*temp)
			continue
		}
		if temp.Dst == "" {
			ga := g.ClientGossiperHandling(*temp)
			go g.rumourMongering(ga.Msg.Rumour, "")
		} else {
			g.SendPrivateMessageFromUser(temp.Dst, temp.Msg)
		}
	}
}
