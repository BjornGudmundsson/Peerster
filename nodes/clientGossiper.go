package nodes

import (
	"encoding/hex"
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
	if route == "/index" || route == "/" {
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

	if route == "/GetRoutingTableJS" {
		fmt.Println("TableJS")
		g.GetRoutingTableJS(wr, req)
		return
	}

	if route == "/PrivateMessage" {
		g.SendPrivateMessageToUser(wr, req)
		return
	}

	if route == "/PostPrivateMessage" {
		g.PostPrivateMessage(wr, req)
		return
	}
	if route == "/AddFile" {
		g.AddFile(wr, req)
		return
	}
	if route == "/RequestFile" {
		g.RequestFile(wr, req)
		return
	}
	if route == "/GetMetaFiles" {
		g.GetMetaFiles(wr, req)
		return
	}
	if route == "/GetFiles" {
		g.GetFiles(wr, req)
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
	if len(text) == 0 || text == "" {
		return
	}
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
	tpl.ExecuteTemplate(wr, "routing.gohtml", g.RoutingTable.Table)
}

//GetRoutingTableJS serves the js file associated with the RoutingTable
//route page.
func (g *Gossiper) GetRoutingTableJS(wr http.ResponseWriter, req *http.Request) {
	http.ServeFile(wr, req, "scripts/routingTable.js")
}

type privateMsg struct {
	Name     string
	Messages []string
}

const hoplimit uint32 = 10

//SendPrivateMessageToUser send the html corresponding to opening a chat dialogue with
//a known peer
func (g *Gossiper) SendPrivateMessageToUser(wr http.ResponseWriter, req *http.Request) {
	name := req.URL.Query()["name"][0]
	messages := g.PrivateMessageStorage.GetMessagesFromOrigin(name)
	priv := privateMsg{
		Name:     name,
		Messages: messages,
	}
	tpl.ExecuteTemplate(wr, "privateMessages.gohtml", priv)
}

//PostPrivateMessage is a route where the user can post his private message
func (g *Gossiper) PostPrivateMessage(wr http.ResponseWriter, req *http.Request) {
	priv := &data.PrivateMessage{
		Origin:      g.Name,
		ID:          0,
		Text:        req.FormValue("text"),
		Destination: req.FormValue("name"),
		HopLimit:    hoplimit,
	}
	gp := &data.GossipPacket{
		PrivateMessage: priv,
	}
	buf, _ := protobuf.Encode(gp)
	conn, e := net.Dial("udp", g.address.String())
	defer conn.Close()
	if e != nil {
		log.Fatal(e)
	}
	conn.Write(buf)
}

func (g *Gossiper) AddFile(wr http.ResponseWriter, req *http.Request) {
	file, fileHeader, e := req.FormFile("file")
	if e != nil {
		log.Fatal(e)
	}
	g.HandleNewFile(fileHeader, file)
	fmt.Println(g.Chunks)
}

func (g *Gossiper) RequestFile(wr http.ResponseWriter, req *http.Request) {
	fn := req.FormValue("fileName")
	mf := req.FormValue("metafile")
	data, e := hex.DecodeString(mf)
	if e != nil {
		fmt.Println("Not a valid hexstring")
		return
	}
	dst := req.FormValue("destination")
	g.DownLoadAFile(fn, data, dst)
}

func (g *Gossiper) GetMetaFiles(wr http.ResponseWriter, req *http.Request) {
	metafiles := g.Files
	tpl.ExecuteTemplate(wr, "metafiles.gohtml", metafiles)
}

func (g *Gossiper) GetFiles(wr http.ResponseWriter, req *http.Request) {
	chunks := g.Chunks
	tpl.ExecuteTemplate(wr, "TextOfFiles.gohtml", chunks)
}
