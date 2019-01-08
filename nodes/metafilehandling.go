package nodes

import (
	"encoding/hex"
	"errors"
	"log"
	"math/big"

	"github.com/BjornGudmundsson/Peerster/data/hashtable"
)

//PopulateFromMetafile takes in a metafile and the filename
//and populates the chunk to peer mapping with the nodes from the chordx
func (g *Gossiper) PopulateFromMetafile(mf []byte, fn, metafilehash string) {
	n := len(mf)
	for i := 0; i < n; i = i + 32 {
		j := i + 32
		chunk := mf[i:j]
		hexChunk := hex.EncodeToString(chunk)
		position := hashtable.HashStringInt(hexChunk)
		pos1, pos2 := g.ChordTable.GetPlaceInChord(position)
		positions := []*big.Int{pos1, pos2}
		for _, pos := range positions {
			if pos != nil {
				owner := g.ChordTable.GetNodeAtPosition(pos)
				if owner == "" {
					log.Fatal(errors.New("The owner was not in the chord. Something is up"))
				}
				g.ChunkToPeer.AddOwnerTochunk(metafilehash, uint64(i/32)+1, owner)
			}
		}
	}
	//g.DownloadingFile(fn)
}
