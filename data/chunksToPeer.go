package data

import (
	"fmt"
	"math/rand"
)

//FileMatcher is used to check
//if there is a match for a metafilehash
type FileMatcher struct {
	MetafileHash string
	ChunkCount   uint64
}

//ChunkToPeer maps chunks of a metafile to
//all of the peers that happen to have these chunks.
//The MetaFileToOwner is a hack for simple direct download
type ChunkToPeer struct {
	MetaFileToPeers map[string]map[uint64][]string
	MetaFileToOwner map[string][]string
	OwnersOfFile    map[string][]string
}

//NewChunkToPeer returns a new pointer
//to an empty ChunkToPeer struct
func NewChunkToPeer() *ChunkToPeer {
	return &ChunkToPeer{
		MetaFileToPeers: make(map[string]map[uint64][]string),
		MetaFileToOwner: make(map[string][]string),
		OwnersOfFile:    make(map[string][]string),
	}
}

//GetRandomOwnerOfMetafile returns a random peer that is known to have the metafile
func (ctp *ChunkToPeer) GetRandomOwnerOfMetafile(metafilehash string) string {
	owners, ok := ctp.MetaFileToOwner[metafilehash]
	if !ok {
		return ""
	}
	ran := rand.Int() % len(owners)
	return owners[ran]
}

//AddOwnerForMetafileHash registers an owner for a particular metafilehash
func (ctp *ChunkToPeer) AddOwnerForMetafileHash(owner, metafilehash string) {
	_, ok := ctp.MetaFileToOwner[metafilehash]
	if !ok {
		ctp.MetaFileToOwner[metafilehash] = make([]string, 0)
	}
	owners := ctp.MetaFileToOwner[metafilehash]
	owners = append(owners, owner)
	ctp.MetaFileToOwner[metafilehash] = owners
}

//AddOwnerTochunk adds an owner to a chunk
func (ctp *ChunkToPeer) AddOwnerTochunk(metafilehash string, index uint64, owner string) {
	_, ok := ctp.MetaFileToPeers[metafilehash]
	if !ok {
		ctp.MetaFileToPeers[metafilehash] = make(map[uint64][]string)
	}
	chunks := ctp.MetaFileToPeers[metafilehash]
	owners := chunks[index]
	owners = append(owners, owner)
	chunks[index] = owners
	ctp.MetaFileToPeers[metafilehash] = chunks
}

//SetOwnerOfMetafileHash sets an owner for a metafile.
func (ctp *ChunkToPeer) SetOwnerOfMetafileHash(metafilehash, owner string) {
	ctp.OwnersOfFile[metafilehash] = append(ctp.OwnersOfFile[metafilehash], owner)
}

//GetRandomOwnerOfChunk takes in a metafilehash and index and returns a random peer
//that supposedly has that chunk.
func (ctp *ChunkToPeer) GetRandomOwnerOfChunk(metafilehash string, index uint64) string {
	chunkKeepers, ok := ctp.MetaFileToPeers[metafilehash][index]
	if ok {
		ran := rand.Int() % len(chunkKeepers)
		chunkKeeper := chunkKeepers[ran]
		return chunkKeeper
	}
	//If no one has the chunk, ask the owner. This is a hack to be compatible with HW2
	owners, ok := ctp.OwnersOfFile[metafilehash]
	if !ok {
		fmt.Println("No one has this chunk or is a registered owner of the metafile")
		//Something is up
		return ""
	}
	ran := rand.Int() % len(owners)
	return owners[ran]
}

//FileMatches checks how many of these filematches are full matches
func (ctp *ChunkToPeer) FileMatches(filematches []FileMatcher) uint64 {
	var temp uint64
	for _, filematch := range filematches {
		chunks := ctp.MetaFileToPeers[filematch.MetafileHash]
		count := filematch.ChunkCount
		match := true
		for i := uint64(0); i < count; i++ {
			if len(chunks[i]) == 0 {
				match = false
				break
			}
		}
		if match {
			temp++
		}
	}
	return temp
}
