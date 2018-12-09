package data

import (
	"encoding/hex"
	"strings"
)

//FoundFile  holds all the
//information of a found file.
type FoundFile struct {
	Origin       string
	ChunkMap     []uint64
	ChunkSize    uint64
	MetafileHash []byte
}

//NewFoundFile returns a new FoundFile structure
func NewFoundFile(src string, result *SearchResult) FoundFile {
	return FoundFile{
		Origin:       src,
		ChunkSize:    result.ChunkCount,
		ChunkMap:     result.ChunkMap,
		MetafileHash: result.MetafileHash,
	}
}

//IsMatch is a function on a FoundFile struct
//that says if this FoundFile is a match for the complete file
func (ff *FoundFile) IsMatch() bool {
	return uint64(len(ff.ChunkMap)) == ff.ChunkSize
}

//FoundFileRepository is a repository for foundfile
//objects and their corresponding filenames
type FoundFileRepository map[string][]FoundFile

//NewFoundFileRepository is a function that returns a new
//FoundFileRepository
func NewFoundFileRepository() FoundFileRepository {
	return FoundFileRepository(make(map[string][]FoundFile))
}

//AddSearchReply adds a search reply to the repository.
func (ffr FoundFileRepository) AddSearchReply(result *SearchResult, src string) {
	foundfile := NewFoundFile(src, result)
	ffr[result.FileName] = append(ffr[result.FileName], foundfile)
}

//FindFullMatches find all of the full matches that have been found that
//match any of the listed keywords and are a full match
func (ffr FoundFileRepository) FindFullMatches(keywords []string) []FileMatcher {
	var fullMatches []FileMatcher
	for filename, matches := range ffr {
		var contains bool
		for _, keyword := range keywords {
			if strings.Contains(filename, keyword) {
				contains = true
				break
			}
		}
		if !contains {
			continue
		}
		if len(matches) == 0 {
			continue
		}
		hx := hex.EncodeToString(matches[0].MetafileHash)
		fm := FileMatcher{
			MetafileHash: hx,
			ChunkCount:   matches[0].ChunkSize,
		}
		fullMatches = append(fullMatches, fm)
	}
	return fullMatches
}
