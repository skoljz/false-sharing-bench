package main_test

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

type WithoutPadding struct {
	a int64
	b int64
}

type WithPadding struct {
	a int64
	_ [CacheLineSize]byte
	b int64
}

const CacheLineSize = 64

func BenchmarkFalseSharing(b *testing.B) {
	runtime.GOMAXPROCS(2)
	s := &WithoutPadding{a: 0, b: 0}
	var wg sync.WaitGroup
	b.ResetTimer()

	wg.Add(2)
	go func() {
		defer wg.Done()
		for range b.N {
			atomic.AddInt64(&s.a, 1)
		}
	}()

	go func() {
		defer wg.Done()
		for range b.N {
			atomic.AddInt64(&s.b, 1)
		}
	}()

	wg.Wait()
}

func BenchmarkWithPadding(b *testing.B) {
	runtime.GOMAXPROCS(2)
	s := &WithPadding{a: 0, b: 0}
	var wg sync.WaitGroup
	b.ResetTimer()

	wg.Add(2)
	go func() {
		defer wg.Done()
		for range b.N {
			atomic.AddInt64(&s.a, 1)
		}
	}()

	go func() {
		defer wg.Done()
		for range b.N {
			atomic.AddInt64(&s.b, 1)
		}
	}()

	wg.Wait()
}
