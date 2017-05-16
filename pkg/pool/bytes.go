// Package pool provides pool structures to help reduce garbage collector pressure.
package pool

import (
	"sync"

	"github.com/CAFxX/gcnotifier"
)

type bytes struct {
	sync.Mutex
	pool [][]byte
}

// Bytes is a pool of byte slices that can be re-used.  Slices in this pool will
// be opportunistically GCed.
type Bytes struct {
	b   *bytes
	gcn *gcnotifier.GCNotifier
}

// NewBytes returns a Bytes pool.
func NewBytes(_ int) *Bytes {
	p := &Bytes{b: &bytes{}, gcn: gcnotifier.New()}
	go p.b.collector(p.gcn.AfterGC())
	return p
}

// Get returns a byte slice size with at least sz capacity. Items
// returned may not be in the zero state and should be reset by the
// caller.
func (p *Bytes) Get(sz int) []byte {
	return p.b.get(sz)
}

// Put returns a slice back to the pool.
func (p *Bytes) Put(c []byte) {
	p.b.put(c)
}

func (b *bytes) get(sz int) []byte {
	b.Lock()
	if pl := len(b.pool); pl > 0 {
		buf := b.pool[pl-1]
		if cap(buf) >= sz {
			b.pool[pl-1] = nil
			b.pool = b.pool[:pl-1]
			b.Unlock()
			return buf[:sz]
		}
	}
	b.Unlock()
	return make([]byte, sz)
}

func (b *bytes) put(c []byte) {
	b.Lock()
	b.pool = append(b.pool, c)
	b.Unlock()
}

func (b *bytes) collector(afterGC <-chan struct{}) {
	for range afterGC {
		b.Lock()
		b.pool = nil
		b.Unlock()
	}
}
