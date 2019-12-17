package gateway

import (
	"bytes"
	"sync"
)

// Buffer wrapper for bytes.Buffer with Free method
type Buffer struct {
	*bytes.Buffer
}

// Free resets buffer and store it to the pool
func (b *Buffer) Free() {
	b.Reset()
	bufPool.Put(b)
}

// NewBuffer create new wrapped buffer
func NewBuffer() *Buffer {
	return bufPool.Get().(*Buffer)
}

var bufPool = sync.Pool{
	New: func() interface{} {
		return &Buffer{new(bytes.Buffer)}
	},
}
