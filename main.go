package main

import (
	"flag"

	"github.com/BjornGudmundsson/Peerster/data"
	"github.com/BjornGudmundsson/Peerster/nodes"
)

func main() {
	port := flag.Int("UIPort", 8080, "The port number")
	name := flag.String("name", "nodeA", "The name of the node")
	addr := flag.String("gossipAddr", "127.0.0.1", "The home address")
	peers := flag.String("peers", "127.0.0.1", "The list of peers")
	simple := flag.Bool("simple", true, "Is it simple")
	rtTimer := flag.Int("rttimer", 0, "timer of rumour chat messages")
	if *simple {
	}
	flag.Parse()
	fp := data.FormatPeers(*peers)
	g := nodes.NewGossiper(*addr, *name, fp, *port)
	go g.ReceiveMessages()
	go g.ClientMessageReceived(*port)
	go g.TCPServer(8080)
	go g.RumourChatting(*rtTimer)
	go g.AntiEntropy()
	g.MiningThread()

}
