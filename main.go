package main

import (
	"flag"
	"fmt"

	"github.com/BjornGudmundsson/Peerster/data"
	"github.com/BjornGudmundsson/Peerster/nodes"
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
	fmt.Println(port)
	fp := data.FormatPeers(*peers)
	g := nodes.NewGossiper(*addr, *name, fp)
	go g.ReceiveMessages()
	go g.ClientMessageReceived(*port)
	for {

	}

}
