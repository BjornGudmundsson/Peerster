package main

import (
	"CS-438/data"
	"CS-438/nodes"
	"flag"
	"fmt"
	"strconv"
	"strings"
)

func main() {
	port := flag.Int("UIPort", 8080, "The port number")
	name := flag.String("name", "nodeA", "The name of the node")
	addr := flag.String("gossipAddr", "127.0.0.1", "The home address")
	peers := flag.String("peers", "127.0.0.1", "The list of peers")
	simple := flag.Bool("simple", true, "Is is simple")
	if *simple {
	}
	flag.Parse()
	Name := *name
	Addr := *addr
	sp := fmt.Sprintf(":%d", *port)
	fp := data.FormatPeers(*peers)
	mainNode := &nodes.Node{
		Messages:   make([]string, 0),
		Neighbours: fp,
		Addr:       Addr,
		Port:       sp,
		Name:       Name,
	}
	go mainNode.EstablishUDPServer()
	fmt.Println(fp)
	for j, por := range fp {
		i := strings.Index(por, ":")
		p := por[i:len(por)]
		fmt.Println(addr)
		newMessages := make([]string, 0)
		s := strconv.FormatInt(int64(j), 10)
		nNode := &nodes.Node{
			Addr:       por,
			Messages:   newMessages,
			Port:       p,
			Name:       "Node " + s,
			Neighbours: fp,
		}
		fmt.Println(p, "bjorn")
		go nNode.EstablishUDPServer()
	}
	//Hardcoded at first, don't know how they want to handle this
	go mainNode.HandleClientConnection()
	for {
	}
}
