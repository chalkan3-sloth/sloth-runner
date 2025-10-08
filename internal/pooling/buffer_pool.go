package pooling

import (
	"bytes"
	"sync"
)

// BufferPool provides a pool of reusable byte buffers to reduce GC pressure
type BufferPool struct {
	pool sync.Pool
}

// NewBufferPool creates a new buffer pool
func NewBufferPool() *BufferPool {
	return &BufferPool{
		pool: sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}
}

// Get retrieves a buffer from the pool
func (p *BufferPool) Get() *bytes.Buffer {
	buf := p.pool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

// Put returns a buffer to the pool
func (p *BufferPool) Put(buf *bytes.Buffer) {
	// Don't pool buffers that are too large (> 64KB)
	if buf.Cap() > 64*1024 {
		return
	}
	buf.Reset()
	p.pool.Put(buf)
}

// ByteSlicePool provides a pool of reusable byte slices
type ByteSlicePool struct {
	pool sync.Pool
	size int
}

// NewByteSlicePool creates a new byte slice pool
func NewByteSlicePool(size int) *ByteSlicePool {
	return &ByteSlicePool{
		size: size,
		pool: sync.Pool{
			New: func() interface{} {
				b := make([]byte, size)
				return &b
			},
		},
	}
}

// Get retrieves a byte slice from the pool
func (p *ByteSlicePool) Get() []byte {
	b := p.pool.Get().(*[]byte)
	return (*b)[:p.size]
}

// Put returns a byte slice to the pool
func (p *ByteSlicePool) Put(b []byte) {
	p.pool.Put(&b)
}
