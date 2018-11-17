package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/BjornGudmundsson/Peerster/data"
	"github.com/dedis/protobuf"
)

func main() {
	UIPort := flag.Int("UIPort", 8080, "This is the port for the user client")
	msg := flag.String("msg", "Hello", "A message to be sent")
	ip := flag.String("ip", "127.0.0.1", "default IP address ")
	dst := flag.String("dest", "", "The destination node of a message or a file download")
	file := flag.String("file", "", "The file to be indexed by the gossiper, or filename of the requested file")
	req := flag.String("request", "", "request a chunk or metafile of this hash")
	keywords := flag.String("keywords", "", "The keywords to search for a file by")
	budget := flag.Int("budget", 2, "The budget for a request")
	flag.Parse()
	sendMessage(*UIPort, *msg, *ip, *dst, *file, *req, *keywords, uint64(*budget))
}

func sendMessage(port int, msg string, ip string, dst string, file, req string, kw string, b uint64) {
	s := fmt.Sprintf("%v:%v", ip, port)
	fmt.Println(s)
	tmsg := &data.TextMessage{
		Msg:      msg,
		Dst:      dst,
		File:     file,
		Request:  req,
		Keywords: kw,
		Budget:   b,
	}
	fmt.Println("Text message: ", *tmsg)
	buf, _ := protobuf.Encode(tmsg)
	conn, e := net.Dial("udp", s)
	defer conn.Close()
	if e != nil {
		log.Fatal(e)
	}
	conn.Write(buf)

}
