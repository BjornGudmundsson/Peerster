package nodes

import (
	"fmt"
	"net/http"
)

//Starts an HTTP server for the gossiper
//client on the specified port
func (g *Gossiper) TCPServer(port int) {
	portStr := fmt.Sprintf(":%v", port)
	http.Handle("/GetMessages", g)
	http.ListenAndServe(portStr, g)
}

func (g *Gossiper) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	var s string
	messages := g.Messages.Messages
	for key, val := range messages {
		for _, rm := range val {
			str := fmt.Sprintf(" Origin %v ID %v COntent %v ", key, rm.ID, rm.Text)
			fmt.Println(str)
			s += str
		}
		s += "\n"
	}
	fmt.Println("yo", s)
	wr.Write([]byte(s))
	g.Messages.PrintMessages()
}
