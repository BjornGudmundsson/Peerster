package data

//MetaData is a struct that
//keeps all of the necessary
//information about a file.
type MetaData struct {
	FileName       string
	FileSize       uint64
	MetaFile       []byte
	HashOfMetaFile string
}

//DataRequest is a struct that
//keeps track of all the information
//necessary to request information from
//a file.
type DataRequest struct {
	Origin      string
	Destination string
	HopLimit    uint32
	HashValue   []byte
}

//DataReply is a struct
//that keeps track of all
//the data necessary to reply
//to a DataRequest
type DataReply struct {
	Origin      string
	Destination string
	HopLimit    uint32
	HashValue   []byte
	Data        []byte
}
