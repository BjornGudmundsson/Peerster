package data

import (
	"sync"
)

//FileFinder is a struct that
//keeps track of where a file
//with a specific file is located
type FileFinder struct {
	MetaFile []byte
	Origin   string
}

//NumberOfChunks is a function that returns the
//number of chunks that a corresponding file has.
func (ff *FileFinder) NumberOfChunks() uint64 {
	return uint64(len(ff.MetaFile))
}

type StateFileFinder struct {
	mux   sync.Mutex
	State map[string]FileFinder
}

//func (sff *StateFileFinder) CheckForKeyWords(keywords []string)
