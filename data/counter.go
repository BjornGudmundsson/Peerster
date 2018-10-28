package data

import (
	"sync"
)

//Counter is a struct that allows concurrent counting methods
type Counter struct {
	c   uint32
	mux sync.Mutex
}

//IncrementAndReturn increments the value for the counter in a
//concurrent way and returns the new value
func (c *Counter) IncrementAndReturn() uint32 {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.c = c.c + 1
	return c.c
}

//ReturnCounter returns the current value of the counter
func (c *Counter) ReturnCounter() uint32 {
	return c.c
}
