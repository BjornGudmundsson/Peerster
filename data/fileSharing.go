package data

//StateDataRequests is a map that maps
//keywords to file structs.
type File struct {
	Name         string
	Size         int64
	MetafileHash []byte
}
