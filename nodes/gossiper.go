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
	address    *net.UDPAddr
	conn       *net.UDPConn
	Name       string
	Neighbours *data.Neighbours
	Messages   *data.MessageHolder
	Counter    *data.Counter
	Status     map[string]chan *data.GossipPacket
	Mongering  *Mongerers
}

//NewGossiper is a function that returns a pointer
//to a new instance of a Gossiper with an active
//UDP connection to the specified address.
func NewGossiper(address, name string, neighbours []string) *Gossiper {
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
	counter := &data.Counter{}
	return &Gossiper{
		Name:       name,
		address:    udpaddr,
		conn:       udpconn,
		Neighbours: conN,
		Messages:   data.NewMessageHolder(),
		Counter:    counter,
		Status:     make(map[string]chan *data.GossipPacket),
		Mongering:  &Mongerers{},
	}
}

//SendRumourMessagesToNeighbours is a function that is bound to a pointer for a gossiper
//It takes a slice of strings as an argument where every string is of a
//conventional IPv4 IP address format like address:port
func (g *Gossiper) SendRumourMessagesToNeighbours(msg string) {
	id := g.Counter.IncrementAndReturn()
	for ip := range g.Neighbours.Neighbours {
		fmt.Println("The id is ", id)
		rm := &data.RumourMessage{
			Origin: g.Name,
			ID:     id,
			Text:   msg,
		}
		gp := &data.GossipPacket{Rumour: rm}
		g.sendRumourMessageToNeighbour(gp, ip)
	}
}

func (g *Gossiper) sendRumourMessageToNeighbour(msg *data.GossipPacket, addr string) {
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
	gp := &data.GossipPacket{}
	gossipChannel := make(chan *GossipAddress)
	go g.delegateMessages(gossipChannel)
	for {
		_, addr, _ := conn.ReadFromUDP(buffer[:])
		protobuf.Decode(buffer, gp)
		go g.Neighbours.AddANeighbour(addr.String())
		gAddress := &GossipAddress{
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
		n, _, e := conn.ReadFromUDP(packet)
		if e != nil {
			fmt.Println(n, e.Error())
		}
		protobuf.Decode(packet, temp)
		fmt.Printf("CLIENT MESSAGE: %v", temp.Msg)
		go g.SendRumourMessagesToNeighbours(temp.Msg)
	}
}
