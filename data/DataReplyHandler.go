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
	currentChunks [][]byte
	finished      bool
	currentIndex  int
	chunks        []ChunkAndMessage
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

type ChunkAndMessage struct {
	Message string
	Chunk   []byte
}

func (dr *DataReplyHandler) Start(fn string, nxtChunk string) {
	dr.isPending = true
	dr.Name = fn
	dr.nextChunk = nxtChunk
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

func (dr *DataReplyHandler) GetCurrentChunks() [][]byte {
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
		dr.nextChunk = hex.EncodeToString(data[i : i+32])
		nxtChunk, _ := hex.DecodeString(dr.nextChunk)
		fmt.Println("First chunk", nxtChunk)
		return nxtChunk, false, true, true
	}
	j := i + 32
	nxtChunk := dr.MetaFile[i:j]
	if hashHex != dr.nextChunk {
		fmt.Println("Next chunk was after all", hex.EncodeToString(nxtChunk))
		return nxtChunk, false, false, false
	}
	dr.currentChunks = append(dr.currentChunks, hashvalue)
	if j == len(dr.MetaFile) {
		fmt.Println("Download completed, got all da chunks")
		dr.finished = true
		return nxtChunk, true, false, false
	}
	dr.nextChunk = hex.EncodeToString(dr.MetaFile[j : j+32])
	dr.currentIndex = j
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

func (dr *DataReplyHandler) GetChunks() []ChunkAndMessage {
	return dr.chunks
}
