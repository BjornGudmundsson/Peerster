package data

import (
	"encoding/hex"
)

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
