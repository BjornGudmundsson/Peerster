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
	Neighbours []string
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
	return &Gossiper{
		Name:       name,
		address:    udpaddr,
		conn:       udpconn,
		Neighbours: neighbours,
	}
}

//SendMessagesToNeighbours is a function that is bound to a pointer for a gossiper
//It takes a slice of strings as an argument where every string is of a
//conventional IPv4 IP address format like address:port
func (g *Gossiper) SendMessagesToNeighbours(msg string) {
	for _, ip := range g.Neighbours {
		udpaddr, e := net.ResolveUDPAddr("udp4", ip)
		if e != nil {
			log.Fatal(e)
		}
		fmt.Println(udpaddr.Port)
		sm := data.NewSimpleMessage(g.Name, msg, ip)
		packetBytes, _ := protobuf.Encode(sm)
		g.conn.WriteToUDP(packetBytes, udpaddr)
	}
}

//This function listens for incoming messages coming
//in on the UDP connection for this node.
func (g *Gossiper) ReceiveMessages() {
	buffer := make([]byte, 1024)
	conn := g.conn
	smp := &data.SimpleMessage{}
	for {
		conn.ReadFromUDP(buffer[:])
		protobuf.Decode(buffer, smp)
		fmt.Println((*smp).Contents)
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
		fmt.Println(udpAddr.IP.String(), udpAddr.Port)
		n, _, e := conn.ReadFromUDP(packet)
		if e != nil {
			fmt.Println(n, e.Error())
		}
		protobuf.Decode(packet, temp)
		fmt.Println(temp.Msg)
		g.SendMessagesToNeighbours(temp.Msg)
	}
}
