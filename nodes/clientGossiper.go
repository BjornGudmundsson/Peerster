package nodes

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"

	"github.com/BjornGudmundsson/Peerster/data"
	"github.com/dedis/protobuf"
)

//Starts an HTTP server for the gossiper
//client on the specified port
func (g *Gossiper) TCPServer(port int) {
	portStr := fmt.Sprintf(":%v", port)
	http.Handle("/GetMessages", g)
	http.ListenAndServe(portStr, g)
}

func (g *Gossiper) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	tpl = template.Must(template.ParseGlob("./templates/*.gohtml"))
	route := req.URL.Path
	if route == "/index" {
		g.GetIndexPage(wr, req)
		return
	}
	if route == "/GetIndexJS" {
		g.GetIndexJS(wr, req)
		return
	}
	if route == "/GetMessages" {
		g.GetMessages(wr, req)
		return
	}
	if route == "/AddMessage" {
		g.AddMessage(wr, req)
		return
	}
	if route == "/GetRoutingTable" {
		g.GetRoutingTable(wr, req)
		return
	}
}

var tpl *template.Template

type tplVars struct {
	Name     string
	Messages []data.RumourMessage
}

//GetIndexPage sends an html page back as a response with the gossiper information
func (g *Gossiper) GetIndexPage(wr http.ResponseWriter, req *http.Request) {
	tv := tplVars{
		Name: g.Name,
	}
	tpl.ExecuteTemplate(wr, "index.gohtml", tv)
}

//GetIndexJS serves the javascript associated with the index webpabe
func (g *Gossiper) GetIndexJS(wr http.ResponseWriter, req *http.Request) {
	http.ServeFile(wr, req, "scripts/index.js")
}

//GetMessages sends
func (g *Gossiper) GetMessages(wr http.ResponseWriter, req *http.Request) {
	//messages := g.Messages.GetMessageString()
	a := g.Messages.GetMessagesInOrder()
	as := tplVars{
		Name:     g.Name,
		Messages: a,
	}
	tpl.ExecuteTemplate(wr, "message.gohtml", as)
}

//AddMessage takes in a message from the user in a form and adds it
func (g *Gossiper) AddMessage(wr http.ResponseWriter, req *http.Request) {
	text := req.FormValue("text")
	addr := g.address
	tmsg := &data.TextMessage{
		Msg: text,
	}
	buf, _ := protobuf.Encode(tmsg)
	cAddr := fmt.Sprintf("%v:%v", addr.IP, g.UIPort)
	conn, e := net.Dial("udp", cAddr)
	defer conn.Close()
	if e != nil {
		log.Fatal(e)
	}
	conn.Write(buf)
}

//GetRoutingTable displays the routing table for debug purposes
func (g *Gossiper) GetRoutingTable(wr http.ResponseWriter, req *http.Request) {
	m := struct{
		Table  map[string]string
	}
	n = m{
		Table : g.RoutingTable.Table,
	}
	tpl.ExecuteTemplate(wr, "routing.gohtml", n)
}
