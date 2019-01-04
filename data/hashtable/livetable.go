package hashtable

import (
	"sync"
	"time"
)

//LiveTable is a struct that keeps track
//of the last time they got a rumour from
//any given node
type LiveTable struct {
	mux       sync.Mutex
	LiveTable map[string]time.Time
}

//NewLiveTable returns a new empty live table
func NewLiveTable() *LiveTable {
	return &LiveTable{
		LiveTable: make(map[string]time.Time),
	}
}

//Update updates the live table with the current time
func (lt *LiveTable) Update(s string) {
	lt.mux.Lock()
	defer lt.mux.Unlock()
	lt.LiveTable[s] = time.Now()
}

//DeadNodes returns the list of dead nodes on the network
func (lt *LiveTable) DeadNodes(d time.Duration) []string {
	lt.mux.Lock()
	defer lt.mux.Unlock()
	temp := make([]string, 0)
	for node, t := range lt.LiveTable {
		if time.Since(t) > d {
			temp = append(temp, node)
		}
	}
	return temp
}
