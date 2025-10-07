package pooling

import (
	"bufio"
	"bytes"
	"sync"
)

// BufferPool manages reusable buffers and scanners for streaming operations
type BufferPool struct {
	bufferPool  sync.Pool
	scannerPool sync.Pool
}

// Global buffer pool for streaming operations
var GlobalBufferPool = NewBufferPool()

// NewBufferPool creates a new buffer pool
func NewBufferPool() *BufferPool {
	return &BufferPool{
		bufferPool: sync.Pool{
			New: func() interface{} {
				// Pre-allocate 4KB buffers (optimal for most streaming)
				return bytes.NewBuffer(make([]byte, 0, 4096))
			},
		},
		scannerPool: sync.Pool{
			New: func() interface{} {
				return bufio.NewScanner(nil)
			},
		},
	}
}

// GetBuffer gets a buffer from the pool
func (p *BufferPool) GetBuffer() *bytes.Buffer {
	buf := p.bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

// PutBuffer returns a buffer to the pool
func (p *BufferPool) PutBuffer(buf *bytes.Buffer) {
	if buf.Cap() > 64*1024 { // Don't pool buffers > 64KB
		return
	}
	buf.Reset()
	p.bufferPool.Put(buf)
}

// GetScanner gets a scanner from the pool
func (p *BufferPool) GetScanner(r interface{ Read([]byte) (int, error) }) *bufio.Scanner {
	scanner := p.scannerPool.Get().(*bufio.Scanner)

	// Create new buffer for the scanner with optimal size
	buf := make([]byte, 0, 64*1024) // 64KB buffer
	scanner.Buffer(buf, 1024*1024)  // Max token size 1MB

	// Reset the scanner with new reader
	scanner.Split(bufio.ScanLines)

	// Use reflection-free reset
	*scanner = *bufio.NewScanner(r)

	return scanner
}

// PutScanner returns a scanner to the pool
func (p *BufferPool) PutScanner(scanner *bufio.Scanner) {
	// Don't need to reset, will be reset on Get
	p.scannerPool.Put(scanner)
}

// SlicePool manages reusable slices for process lists and other arrays
type SlicePool struct {
	processInfoPool sync.Pool
	stringPool      sync.Pool
}

// GlobalSlicePool for reusable slices
var GlobalSlicePool = NewSlicePool()

// NewSlicePool creates a new slice pool
func NewSlicePool() *SlicePool {
	return &SlicePool{
		processInfoPool: sync.Pool{
			New: func() interface{} {
				// Pre-allocate slice for 50 processes (typical)
				slice := make([]interface{}, 0, 50)
				return &slice
			},
		},
		stringPool: sync.Pool{
			New: func() interface{} {
				// Pre-allocate slice for 20 strings (typical)
				slice := make([]string, 0, 20)
				return &slice
			},
		},
	}
}

// GetProcessSlice gets a process slice from the pool
func (p *SlicePool) GetProcessSlice() *[]interface{} {
	slice := p.processInfoPool.Get().(*[]interface{})
	*slice = (*slice)[:0] // Clear but keep capacity
	return slice
}

// PutProcessSlice returns a process slice to the pool
func (p *SlicePool) PutProcessSlice(slice *[]interface{}) {
	if cap(*slice) > 200 { // Don't pool huge slices
		return
	}
	p.processInfoPool.Put(slice)
}

// GetStringSlice gets a string slice from the pool
func (p *SlicePool) GetStringSlice() *[]string {
	slice := p.stringPool.Get().(*[]string)
	*slice = (*slice)[:0] // Clear but keep capacity
	return slice
}

// PutStringSlice returns a string slice to the pool
func (p *SlicePool) PutStringSlice(slice *[]string) {
	if cap(*slice) > 100 { // Don't pool huge slices
		return
	}
	p.stringPool.Put(slice)
}
