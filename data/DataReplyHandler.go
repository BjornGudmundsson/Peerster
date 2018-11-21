package data

import (
	"encoding/hex"
	"fmt"
)

type DataReplyHandler struct {
	isPending     bool
	Name          string
	nextChunk     string
	MetaFile      []byte
	currentChunks []string
	finished      bool
	currentIndex  int
}

func (dr *DataReplyHandler) Clear() {
	dr.isPending = false
	dr.Name = ""
	dr.nextChunk = ""
	dr.MetaFile = nil
	dr.currentChunks = nil
	dr.finished = false
	dr.currentIndex = 0
}

func NewDataReplyHandler() *DataReplyHandler {
	return &DataReplyHandler{
		isPending:     false,
		Name:          "",
		nextChunk:     "",
		MetaFile:      nil,
		currentChunks: nil,
		finished:      false,
		currentIndex:  0,
	}
}

func (dr *DataReplyHandler) Start(fn string, nxtChunk string, mf []byte, currentIndex int) {
	dr.isPending = true
	dr.Name = fn
	dr.nextChunk = nxtChunk
	dr.MetaFile = mf
	dr.currentIndex = currentIndex
}

func Compare(a []byte, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, val := range a {
		if val != b[i] {
			return false
		}
	}
	return true
}

func (dr *DataReplyHandler) GetCurrentChunks() []string {
	return dr.currentChunks
}

func (dr *DataReplyHandler) Update(hashvalue []byte, data []byte) ([]byte, bool, bool, bool) {
	if !dr.isPending {
		return nil, false, false, false
	}
	i := dr.currentIndex
	hashHex := hex.EncodeToString(hashvalue)
	if hashHex != dr.nextChunk && dr.MetaFile == nil {
		nxtChunk, _ := hex.DecodeString(dr.nextChunk)
		return nxtChunk, false, false, false
	}
	fmt.Println("I did not get an error")
	if dr.MetaFile == nil {
		dr.MetaFile = data
		//Every chunk is exactly 32 bytes
		//return false since I am not finished
		//because I have only received the metafile
		//The chunk i return is the first chunk of the file
		fmt.Println("")
		fmt.Printf("Downloading metafile of %v from anotherPeer", dr.Name)
		dr.nextChunk = hex.EncodeToString(data[i : i+32])
		nxtChunk, _ := hex.DecodeString(dr.nextChunk)
		return nxtChunk, false, true, true
	}
	j := i + 32
	nxtChunk := dr.MetaFile[i:j]
	if hashHex != dr.nextChunk {
		return nxtChunk, false, false, false
	}
	dr.currentChunks = append(dr.currentChunks, hashHex)
	if j == len(dr.MetaFile) {
		dr.finished = true
		return nxtChunk, true, false, false
	}
	dr.nextChunk = hex.EncodeToString(dr.MetaFile[j : j+32])
	dr.currentIndex = j
	fmt.Println("")
	fmt.Printf("Downloading %v chunk %v from anotherPeer", dr.Name, j/32)
	return dr.MetaFile[j : j+32], false, false, true

}

func (dr *DataReplyHandler) IsFinished() bool {
	return dr.finished
}

func (dr *DataReplyHandler) IsPending() bool {
	return dr.isPending
}

func (dr *DataReplyHandler) NextChunk() string {
	return dr.nextChunk
}

//HandlerDataReplies handles data replies. It is a map
//of channels where the keys are the hex values of the
//chunks.
type HandlerDataReplies map[string]chan DataReply

//AddChunk appends a channel to the HandlerDataReplies map with a corresponding chunk
func (hdr HandlerDataReplies) AddChunk(chunk string, ch chan DataReply) {
	hdr[chunk] = ch
}

//AddMetafile takes in a metafile and a channel and adds that channel to the
//HandleDataReplies for every chunk in the metafile
func (hdr HandlerDataReplies) AddMetafile(metafile []byte, ch chan DataReply) {
	n := len(metafile) / 32
	for i := 0; i < n; i++ {
		j := i + 1
		chunk := hex.EncodeToString(metafile[i*32 : j*32])
		hdr.AddChunk(chunk, ch)
	}
}

//DeleteMetafile takes in a metafile and deletes every chunk of
//that metafile from the map
func (hdr HandlerDataReplies) DeleteMetafile(metafile []byte) {
	n := len(metafile) / 32
	for i := 0; i < n; i++ {
		j := i + 1
		chunk := hex.EncodeToString(metafile[i*32 : j*32])
		delete(hdr, chunk)
	}
}

//PassReplyToChannel passes a chunk to the corresponding channelgo
func (hdr HandlerDataReplies) PassReplyToChannel(datareply DataReply) {
	hx := hex.EncodeToString(datareply.HashValue)
	ch, ok := hdr[hx]
	if !ok {
		//The gossiper was not waiting on this chunk
		return
	}
	//Now I have passed the datareply to the channel.
	//It should now be available to the downloading process
	ch <- datareply
}

//NewHandlerDataReplies returns a new empty HandlerDataReplies
func NewHandlerDataReplies() HandlerDataReplies {
	return HandlerDataReplies(make(map[string]chan DataReply, 1024))
}
