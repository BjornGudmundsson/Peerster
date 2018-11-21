package data

import (
	"fmt"
	"strings"
	"sync"
)

//FileFinder is a struct that
//keeps track of where a file
//with a specific file is located
type FileFinder struct {
	MetaFile string
	Origins  []string
}

//NumberOfChunks is a function that returns the
//number of chunks that a corresponding file has.
func (ff *FileFinder) NumberOfChunks() uint64 {
	return uint64(len(ff.MetaFile) / 32)
}

//Check takes in a slice of bytes and checks if
//that slice of bytes corresponds to the fully
//constructed file that this filefinder represents.
func (ff *FileFinder) check(bs []uint64) bool {
	return (len(ff.MetaFile) / 32) == len(bs)
}

//AddOrigin takes in the name of the node where it got a
//SearchReply from and checks if the Chunkmap has all the
//chunks to reconstruct the file and adds that to the list of
//known origins of the file
func (ff *FileFinder) AddOrigin(name string, bs []uint64) bool {
	if !ff.check(bs) {
		return false
	}
	ff.Origins = append(ff.Origins, name)
	return true
}

//StateFileFinder is a struct
//that allows for concurrent access
//to all of the files that the gossiper
//knows exists and where they can be found
type StateFileFinder struct {
	mux        sync.Mutex
	HashToName map[string]string
	State      map[string]FileFinder
}

//NewStateFileFinder returns a new empty
//StateFileFinder.
func NewStateFileFinder() StateFileFinder {
	return StateFileFinder{
		State:      make(map[string]FileFinder),
		HashToName: make(map[string]string),
	}
}

//AddOrigin is a method bound to a StateFileFinder struct. It takes in the filename
//of the file and the name of the node where the file is possibly and tries to add it
//to the corresponding filefinder struct.
func (sff *StateFileFinder) AddOrigin(filename, src string, chunkmap []uint64) {
	filefinder := sff.State[filename]
	ff := &filefinder
	didAdd := ff.AddOrigin(src, chunkmap)
	if didAdd {
		sff.State[filename] = *ff
		fmt.Printf("\nFOUND MATCH for %v at %v \n", filename, src)
	}
}

//WaitingForMetafile is a function on a stateFileFinder that figures out
//if the gossiper is waiting this particular metafilehash and then adds
//the corresponding metafile to a filefinder that corresponds to the filename
//for the given metafilehash.
func (sff *StateFileFinder) WaitingForMetafile(metafilehash, metafile string) (string, bool) {
	name, ok := sff.HashToName[metafilehash]
	if !ok {
		return "", false
	}
	filefinder := FileFinder{
		MetaFile: metafile,
		Origins:  make([]string, 0),
	}
	sff.State[name] = filefinder
	return name, ok
}

//HasMetaFile is a function on a state file finder struct
//that returns a boolean indicating if the Metafile for
//a particular file has been found. The parameters it takes
//are the metafilehash of the file.
func (sff *StateFileFinder) HasMetaFile(mfh string) bool {
	_, ok := sff.HashToName[mfh]
	return ok
}

//HasMatchForFileName is a function to a StateFileFinder that checks if for a
//given filename there is match and how many matches are there for that filename
func (sff *StateFileFinder) HasMatchForFileName(filename string) (bool, int) {
	filefinder, ok := sff.State[filename]
	if !ok {
		return false, 0
	}
	origins := filefinder.Origins
	results := len(origins)
	return results == 0, results
}

//HasMatchForKeyWords is a function bound to a StateFileFinder struct that takes in a set of
//keywords and a threshold. If the number is equal to or less than the threshold indicated in
//the function parameter it stop and returns the amount of mathches it founds. A filename is a match if
//it contains any of the keywords passed to the function.
func (sff *StateFileFinder) HasMatchForKeyWords(keywords []string, threshhold int) (bool, int) {
	sum := 0
	for _, keyword := range keywords {
		for filename := range sff.State {
			hasKeyword := strings.Contains(filename, keyword)
			if !hasKeyword {
				continue
			}
			has, results := sff.HasMatchForFileName(filename)
			if !has {
				continue
			}
			sum = sum + results
			if sum >= threshhold {
				return true, sum
			}
		}
	}
	return false, sum
}
