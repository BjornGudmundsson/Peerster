package data

//MetaFileHashes keeps track of which
//metafiles this Gossiper is waiting
//for.
type MetaFileHashes map[string]string

//Clear is function on a MetaFileHashes object
//that removes a key value pair based on the
//key value of metafilehash
func (mfh MetaFileHashes) Clear(mf string) {
	mfh[mf] = ""
}

//NewMetaFileHashes is just a helper function
//that returns a new empty MetaFileHashes object
func NewMetaFileHashes() MetaFileHashes {
	return MetaFileHashes(make(map[string]string))
}
