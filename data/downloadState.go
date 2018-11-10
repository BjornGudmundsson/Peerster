package data

//FileState is a struct that
//holds onto the MetaFile of a
//file and what chunks it currently has
type FileState struct {
	MetaFile      []byte
	CurrentChunks []string
}

//DownloadState is a map that has the name
//of the file as the key and the state of the
//file as the value.
type DownloadState map[string]FileState

//NewDownloadState returns a new empty
//instance of a DownloadState
func NewDownloadState() DownloadState {
	return DownloadState(make(map[string]FileState))
}
