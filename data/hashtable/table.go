package hashtable

import (
	"math/big"
	"sync"
)

//ChordTable keeps track of who is which
//position on the chord.
//Maps the string representation of
//a big int to an origin
type ChordTable struct {
	mux       sync.Mutex
	table     map[string]string
	positions []*big.Int
}

//NewChordTable returns a new empty
//ChordTable
func NewChordTable() *ChordTable {
	return &ChordTable{
		table:     make(map[string]string),
		positions: make([]*big.Int, 0),
	}
}

//GetNodeAtPosition gets the node at the given position.
func (ct *ChordTable) GetNodeAtPosition(b *big.Int) string {
	sp := b.String()
	return ct.table[sp]
}

//GetTable returns the chord table
func (ct *ChordTable) GetTable() map[string]string {
	return ct.table
}

func (ct *ChordTable) GetPositions() []*big.Int {
	return ct.positions
}

//AddToChord adds the origin to the chord list
func (ct *ChordTable) AddToChord(s string) {
	ct.mux.Lock()
	defer ct.mux.Unlock()
	position := HashStringInt(s)
	sp := position.String()
	_, ok := ct.table[sp]
	if ok {
		return
	}
	ct.table[sp] = s
	ct.addToPositionsList(position)
}

func (ct *ChordTable) addToPositionsList(b *big.Int) {
	n := len(ct.positions)
	temp := make([]*big.Int, 0)
	var inserted bool
	for i := 0; i < n; i++ {
		cmp2 := ct.positions[i].Cmp(b)
		if cmp2 > 0 {
			temp = append(temp, b)
			inserted = true
		}
		temp = append(temp, ct.positions[i])
	}
	if !inserted {
		temp = append(temp, b)
	}
	ct.positions = temp
}

//GetPlaceInChord returns the place in the chord.
func (ct *ChordTable) GetPlaceInChord(b *big.Int) (*big.Int, *big.Int) {
	return FindPlaceInChord(ct.positions, b)
}

//RemoveNode removes a node from the chord table
func (ct *ChordTable) RemoveNode(origin string) {
	ct.mux.Lock()
	defer ct.mux.Unlock()
	position := HashStringInt(origin)
	sp := position.String()
	delete(ct.table, sp)
	ct.removeFromPositionsList(position)
}

func (ct *ChordTable) removeFromPositionsList(b *big.Int) {
	n := len(ct.positions)
	for i := 0; i < n; i++ {
		b2 := ct.positions[i]
		cmp := b2.Cmp(b)
		if cmp == 0 {
			a1 := ct.positions[:i]
			a2 := ct.positions[i+1:]
			ct.positions = append(a1, a2...)
			return
		}
	}
}
