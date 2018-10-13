package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/dedis/protobuf"

	"github.com/BjornGudmundsson/Peerster/data"
)

func main() {
	UIPort := flag.Int("UIPort", 8080, "This is the port for the user client")
	msg := flag.String("msg", "Hello", "A message to be sent")
	ip := flag.String("ip", "127.0.0.1", "default IP address ")
	flag.Parse()
	fmt.Println(*UIPort, *msg)
	sendMessage(*UIPort, *msg, *ip)
}

func sendMessage(port int, msg string, ip string) {
	s := fmt.Sprintf("%v:%v", ip, port)
	/*udpAddr, e := net.ResolveUDPAddr("udp4", s)
	if e != nil {
		log.Fatal(e)
	}*/
	fmt.Println(s)
	tmsg := &data.TextMessage{
		Msg: msg,
	}
	buf, _ := protobuf.Encode(tmsg)
	conn, e := net.Dial("udp", s)
	defer conn.Close()
	if e != nil {
		log.Fatal(e)
	}
	conn.Write(buf)

}
