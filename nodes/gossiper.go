package nodes

import (
	"fmt"
	"log"
	"net"

	"github.com/BjornGudmundsson/Peerster/data"
	"github.com/dedis/protobuf"
)

//Gossiper is an instance of a node running the peerster protocol.
//It has the UDP address of it self and its active UDP connection.
type Gossiper struct {
	address      *net.UDPAddr
	conn         *net.UDPConn
	Name         string
	Neighbours   *data.Neighbours
	Messages     *data.MessageHolder
	Counter      *data.Counter
	Status       *Status
	Mongering    MongererMessages
	enPeer       *EntropyPeer
	UIPort       int
	RoutingTable *data.RoutingTable
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
		Table : make(map[string]string)
	}
	return &Gossiper{
		Name:         name,
		address:      udpaddr,
		conn:         udpconn,
		Neighbours:   conN,
		Messages:     data.NewMessageHolder(),
		Counter:      counter,
		Status:       status,
		Mongering:    mong,
		enPeer:       &EntropyPeer{},
		UIPort:       p,
		RoutingTable: routingTable,
	}
}

func (g *Gossiper) sendMessageToNeighbour(msg *data.GossipPacket, addr string) {
	udpaddr, e := net.ResolveUDPAddr("udp4", addr)
	if e != nil {
		log.Fatal(e)
	}
	packetbyte, _ := protobuf.Encode(msg)
	g.conn.WriteToUDP(packetbyte, udpaddr)
}

//ReceiveMessages listens for incoming messages coming
//in on the UDP connection for this node.
func (g *Gossiper) ReceiveMessages() {
	buffer := make([]byte, 1024)
	conn := g.conn
	gossipChannel := make(chan GossipAddress)
	go g.delegateMessages(gossipChannel)
	for {
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
		id := g.Counter.IncrementAndReturn()
		n, _, e := conn.ReadFromUDP(packet)
		if e != nil {
			fmt.Println(n, e.Error())
		}
		protobuf.Decode(packet, temp)
		fmt.Printf("CLIENT MESSAGE: %v", temp.Msg)
		rm := &data.RumourMessage{
			Origin: g.Name,
			ID:     id,
			Text:   temp.Msg,
		}
		g.Messages.AddAMessage(*rm)
		gp := &data.GossipPacket{
			Rumour: rm,
		}
		ga := &GossipAddress{
			Addr: g.address.String(),
			Msg:  gp,
		}
		go g.rumourMongering(ga)
	}
}
