package hashtable

import "sync"

//ReplyHandler handles the
//Chunk store replies and
//passes them to the appropriate process
type ReplyHandler struct {
	mux       sync.Mutex
	processes map[string]chan *ChunkStoreReply
}

//NewReplyHandler returns a new empty reply handler
func NewReplyHandler() *ReplyHandler {
	return &ReplyHandler{
		processes: make(map[string]chan *ChunkStoreReply),
	}
}

//PassToProcess take a ChunkStoreReply and passes it to the
//corresponding process.
func (rh *ReplyHandler) PassToProcess(csr *ChunkStoreReply) {
	identifier := csr.Src + "-" + csr.Hash
	process, ok := rh.processes[identifier]
	if !ok {
		return
	}
	if process == nil {
		return
	}
	process <- csr
}

//AddProcess add a process to the struct
func (rh *ReplyHandler) AddProcess(identifier string, ch chan *ChunkStoreReply) {
	rh.mux.Lock()
	defer rh.mux.Unlock()
	rh.processes[identifier] = ch
}

//RemoveProcess removes the process with the corresponding identifier.
func (rh *ReplyHandler) RemoveProcess(identifier string) {
	rh.mux.Lock()
	defer rh.mux.Unlock()
	delete(rh.processes, identifier)
}
