package data

import (
	"time"
)

//SearchRequest is a struct holding
//all the necessary variables to perform a
//search for a file over the peerster
//network.
type SearchRequest struct {
	Origin   string
	Budget   uint64
	Keywords []string
}

//RecentRequests is a map that maps a string
//of the form Origin:Keywords to the time
//when it was registered.
type RecentRequests map[string]time.Time

//DeleteByTime deletes the entries where the time
//that has elapsed since it was added to the map
func (rr RecentRequests) DeleteByTime(t time.Time) {
	for key, val := range rr {
		if 10*val.Second() > 10*t.Second()-5 {
			delete(rr, key)
		}
	}
}
