package nodes

import (
	"encoding/hex"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"

	"github.com/BjornGudmundsson/Peerster/data/hashtable"

	"github.com/BjornGudmundsson/Peerster/data/peersterCrypto"

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
	if route == "/GetFoundFiles" {
		g.GetFoundFiles(wr, req)
		return
	}
	if route == "/FoundFile" {
		g.GetFoundFile(wr, req)
		return
	}
	if route == "/GetFoundFileJS" {
		g.GetFoundFilesJS(wr, req)
		return
	}
	if route == "/DownloadFoundFile" {
		g.DownloadFoundFile(wr, req)
		return
	}
	if route == "/DownloadMetaFile" {
		//g.DownloadMetaFile(wr, req)
		return
	}
	if route == "/GetChord" {
		g.GetChordTable(wr, req)
		return
	}
	if route == "/GetAllPublicKeys" {
		g.GetAllPublicKeys(wr, req)
		return
	}
	if route == "/ShareSecretWithPeer" {
		g.ShareSecretWithPeer(wr, req)
		return
	}
	if route == "/GetSecrets" {
		g.GetSecrets(wr, req)
		return
	}
	if route == "/DownloadSecretFile" {
		g.DownloadSecretFile(wr, req)
		return
	}
	if route == "/DownloadFileFromNetwork" {
		g.DownloadFileFromNetwork(wr, req)
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
	a := g.RumourHolder.GetMessagesInOrder()
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

func (g *Gossiper) GetFoundFilesJS(wr http.ResponseWriter, req *http.Request) {
	http.ServeFile(wr, req, "scripts/foundFiles.js")
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
}

func (g *Gossiper) RequestFile(wr http.ResponseWriter, req *http.Request) {
	fn := req.FormValue("fileName")
	mf := req.FormValue("metafile")
	md := data.MetaData{
		FileName:       fn,
		HashOfMetaFile: mf,
		MetaFile:       nil,
		FileSize:       0,
	}
	g.Files[fn] = md
	dst := req.FormValue("destination")
	if dst != "" {
		g.ChunkToPeer.SetOwnerOfMetafileHash(mf, dst)
		g.ChunkToPeer.AddOwnerForMetafileHash(dst, mf)
	}
	go g.DownloadingFile(fn)
}

func (g *Gossiper) DownloadFoundFile(wr http.ResponseWriter, req *http.Request) {
	fn := req.FormValue("name")
	go g.DownloadingFile(fn)
}

func (g *Gossiper) GetMetaFiles(wr http.ResponseWriter, req *http.Request) {
	metafiles := g.Files
	m := make(map[string]string)
	for key, val := range metafiles {
		m[key] = hex.EncodeToString(val.MetaFile)
	}
	tpl.ExecuteTemplate(wr, "metafiles.gohtml", m)
}

func (g *Gossiper) GetFiles(wr http.ResponseWriter, req *http.Request) {
	chunks := g.Chunks
	tpl.ExecuteTemplate(wr, "TextOfFiles.gohtml", chunks)
}

func (g *Gossiper) GetFoundFiles(wr http.ResponseWriter, req *http.Request) {
	foundfiles := g.FoundFileRepository
	tpl.ExecuteTemplate(wr, "FoundFiles.gohtml", foundfiles)
}

type tempStruct struct {
	Name    string
	Matches []data.FoundFile
}

func (g *Gossiper) GetFoundFile(wr http.ResponseWriter, req *http.Request) {
	name := req.URL.Query()["name"][0]
	matches := g.FoundFileRepository[name]
	ts := tempStruct{
		Name:    name,
		Matches: matches,
	}
	tpl.ExecuteTemplate(wr, "FoundFile.gohtml", ts)
}

func (g *Gossiper) GetChordTable(wr http.ResponseWriter, req *http.Request) {
	positions := g.ChordTable.GetPositions()
	m := make([]string, 0)
	for _, val := range positions {
		m = append(m, val.String())
	}
	table := g.ChordTable.GetTable()
	fmt.Println(m)
	tpl.ExecuteTemplate(wr, "chord.gohtml", table)
}

//DownloadMetaFile downloads the metafile with the corresponding file name from
//an HTML form
/*func (g *Gossiper) DownloadMetaFile(wr http.ResponseWriter, req *http.Request) {
	metafile := req.FormValue("metafile")
	metafiledata, e := hex.DecodeString(metafile)
	if e != nil {
		log.Fatal(e)
	}
	fn := req.FormValue("filename")
	if fn == "" {
		log.Fatal(errors.New("Got an empty string from the form"))
	}
	g.PopulateFromMetafile(metafiledata, fn)
}*/

//GetAllPublicKeys is a route that displays all the public key pairs that are
//logged on the longest chain
func (g *Gossiper) GetAllPublicKeys(wr http.ResponseWriter, req *http.Request) {
	pairs := g.GetAllPublicKeyInLongestChain()
	tpl.ExecuteTemplate(wr, "publicKeys.gohtml", pairs)
}

//ShareSecretWithPeer is a GUI route that shares a secret with a peer. The file has
//to have been indexed before being shared or else nothing will happen.
func (g *Gossiper) ShareSecretWithPeer(wr http.ResponseWriter, req *http.Request) {
	fn := req.FormValue("filename")
	peer := req.FormValue("peer")
	if peer == "" {
		peer = g.Name
	}
	g.ShareSecret(fn, peer)
}

//GetSecrets renders a template showing the names of the files
//that have been shared with this peer.
func (g *Gossiper) GetSecrets(wr http.ResponseWriter, req *http.Request) {
	var name string
	val := req.URL.Query()["name"][0]
	if val == "" {
		name = g.Name
	} else {
		name = val
	}
	fmt.Println("Val: ", val, name)
	secrets := g.GetSecretsForPeer(name)
	verifiedSecrets := make([]peersterCrypto.Secret, 0)
	priv := g.PrivateKey
	for _, eSecret := range secrets {
		secret, e := priv.DecryptSecret(&eSecret)
		if e == nil {
			verifiedSecrets = append(verifiedSecrets, *secret)
		}
	}

	//Indexing the file metadata locally for the gossiper.
	for _, secret := range verifiedSecrets {
		mfh := hex.EncodeToString(secret.MetaFileHash)
		p := hashtable.HashStringInt(mfh)
		p1, p2 := g.ChordTable.GetPlaceInChord(p)
		if p1 != nil {
			node := g.ChordTable.GetNodeAtPosition(p1)
			fmt.Println("owner of metafile: ", node)
			g.ChunkToPeer.AddOwnerForMetafileHash(node, mfh)
		}
		if p2 != nil {
			node := g.ChordTable.GetNodeAtPosition(p2)
			fmt.Println("owner of metafile: ", node)
			g.ChunkToPeer.AddOwnerForMetafileHash(node, mfh)
		}
		fmt.Println("IV", secret.IV)
		md := data.MetaData{
			FileName:       secret.FileName,
			FileSize:       7,
			IV:             secret.IV,
			Key:            secret.Key,
			MetaFile:       nil,
			HashOfMetaFile: mfh,
		}
		g.Files[secret.FileName] = md
	}
	tpl.ExecuteTemplate(wr, "secrets.gohtml", verifiedSecrets)
}

func (g *Gossiper) DownloadSecretFile(wr http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()["filename"]
	if len(query) == 0 {
		log.Fatal(errors.New("Did not put in a query parameter"))
	}
	fn := query[0]
	if _, ok := g.Files[fn]; !ok {
		log.Fatal(errors.New("This file has not been indexed"))
	}
	tpl.ExecuteTemplate(wr, "download.gohtml", fn)
}

func (g *Gossiper) DownloadFileFromNetwork(wr http.ResponseWriter, req *http.Request) {
	fn := req.FormValue("filename")
	if fn == "" {
		log.Fatal(errors.New("Got an empty filename"))
	}
	fmt.Println("tcp call again")
	go g.DownloadingFile(fn)
	fmt.Fprintf(wr, "%v", "Bjorn")
}
