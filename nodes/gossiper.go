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
	conN := &data.Neighbours{Neighbours: m}
	return &Gossiper{
		Name:       name,
		address:    udpaddr,
		conn:       udpconn,
		Neighbours: conN,
	}
}

//SendMessagesToNeighbours is a function that is bound to a pointer for a gossiper
//It takes a slice of strings as an argument where every string is of a
//conventional IPv4 IP address format like address:port
func (g *Gossiper) SendMessagesToNeighbours(msg string) {
	for ip := range g.Neighbours.Neighbours {
		udpaddr, e := net.ResolveUDPAddr("udp4", ip)
		if e != nil {
			log.Fatal(e)
		}
		fmt.Println(udpaddr.Port)
		sm := data.NewSimpleMessage(g.Name, msg, ip)
		gp := &data.GossipPacket{Simple: sm}
		packetBytes, _ := protobuf.Encode(gp)
		g.conn.WriteToUDP(packetBytes, udpaddr)
	}
}

//ReceiveMessages listens for incoming messages coming
//in on the UDP connection for this node.
func (g *Gossiper) ReceiveMessages() {
	buffer := make([]byte, 1024)
	conn := g.conn
	gp := &data.GossipPacket{}
	for {
		conn.ReadFromUDP(buffer[:])
		protobuf.Decode(buffer, gp)
		sm := (*gp).Simple
		go g.handleIncomingMessage(sm)

	}
}

func (g *Gossiper) handleIncomingMessage(sm *data.SimpleMessage) {
	fmt.Printf("SIMPLE MESSAGE origin %v from %v content %v \n", sm.OriginalName, sm.RelayPeerAddr, sm.Contents)
	g.Neighbours.PrintNeighbours()
	g.Neighbours.AddANeighbour(sm.RelayPeerAddr)
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
		fmt.Println(udpAddr.IP.String(), udpAddr.Port)
		n, _, e := conn.ReadFromUDP(packet)
		if e != nil {
			fmt.Println(n, e.Error())
		}
		protobuf.Decode(packet, temp)
		fmt.Printf("CLIENT MESSAGE: %v", temp.Msg)
		g.SendMessagesToNeighbours(temp.Msg)
	}
}
