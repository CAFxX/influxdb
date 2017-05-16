package pool

import (
	_bytes "bytes"
	"runtime"
	"testing"
	"time"
)

var hello = []byte("hello")

func TestBytes(t *testing.T) {
	pool := NewBytes(0)

	buf1 := pool.Get(len(hello))
	buf1 = append(buf1[:0], hello...)
	pool.Put(buf1)

	buf2 := pool.Get(len(hello))
	if _bytes.Compare(buf2, hello) != 0 {
		t.Fatalf("buffer not recycled")
	}
	pool.Put(buf2)

	runtime.GC()
	<-time.After(10 * time.Millisecond)

	buf3 := pool.Get(len(hello))
	if _bytes.Compare(buf3, hello) == 0 {
		t.Fatalf("buffer not collected")
	}
}

func BenchmarkBytesFullReuse(b *testing.B) {
	pool := NewBytes(0)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := pool.Get(256)
			pool.Put(buf)
		}
	})
}

func BenchmarkBytesPartialReuse(b *testing.B) {
	pool := NewBytes(0)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := pool.Get(256)
			pool.Put(buf)
			pool.Get(256)
		}
	})
}
