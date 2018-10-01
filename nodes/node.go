package nodes

import (
	"CS-438/data"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
)

//Node is a struct that keeps track of the needed information to
//identify a process within a node. The port variable tells in which
//socket the packet should go to and the IP variable is the address
type Node struct {
	Messages   []string
	Neighbours []string
	Addr       string
	Port       string
	Name       string
}

//PeersterNode is an interface that describes the behaviour
//associated with a node running the peerster application
type PeersterNode interface {
	SendMessageToNeighbours(msg string)
	EstablishUDPServer()
	HandleClientConnection()
	handleClientRequest(wr http.ResponseWriter, req *http.Request)
}

//SendMessageToNeighbours sends a string to all of
//the neighbours of the specified node.
func (n *Node) SendMessageToNeighbours(msg string) {
	for _, val := range n.Neighbours {
		index := strings.Index(val, ":")
		p := val[index+1 : len(val)]
		pi, _ := strconv.Atoi(p)
		ip := val[0:index]
		bs := data.SplitIP(ip)
		c, _ := net.DialUDP("udp", nil, &net.UDPAddr{IP: bs, Port: pi, Zone: ""})
		c.Write([]byte(msg))
	}
}

//EstablishUDPServer instantiates a new UDP server at the
//gossipAddr for the specified node
func (n *Node) EstablishUDPServer() {
	udpAddr, e := net.ResolveUDPAddr("udp4", n.Addr)
	if e != nil {
		log.Fatal(e)
	}
	udpConn, e := net.ListenUDP("udp4", udpAddr)
	if e != nil {
		log.Fatal(e)
	}
	i := strings.Index(n.Addr, ":")
	fmt.Println("Before udp", i)
	udpPort := n.Addr[i:len(n.Addr)]
	fmt.Println(udpPort)
	ln, e := net.ListenPacket("udp", udpPort)
	if e != nil {
		log.Fatal(e)
	}
	defer ln.Close()
	buf := make([]byte, 1024)
	for {
		n, addr, _ := ln.ReadFrom(buf)
		fmt.Println(addr.String())
		fmt.Println("New node")
		fmt.Println(udpAddr.Port)
		val := string(buf[0:n])
		fmt.Println("Message is ", val)
	}
}

//HandleClientConnection handles the connection of
//the client to the peerster system on the device
//that has the specified node
func (n *Node) HandleClientConnection() {
	http.HandleFunc("/postMessage", n.handleClientRequest)
	http.ListenAndServe(n.Port, nil)
}

func (n *Node) handleClientRequest(wr http.ResponseWriter, req *http.Request) {
	msg := req.FormValue("msg")
	n.Messages = append(n.Messages, msg)
	n.SendMessageToNeighbours(msg)
	wr.Write([]byte("route works"))
}
