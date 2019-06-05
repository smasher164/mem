package mem_test

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"unsafe"

	"github.com/smasher164/mem"
)

var szpage = os.Getpagesize()

type slice struct {
	addr unsafe.Pointer
	len  int
	cap  int
}

func errpanic(f func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	f()
	return
}

func TestZero(t *testing.T) {
	p := mem.Alloc(0)
	if p != nil {
		t.Error("mem.Alloc(0) must return nil")
	}
	err := errpanic(func() { mem.Free(p) })
	if err != nil {
		t.Errorf("mem.Free(mem.Alloc(0)) panics: %v", err)
	}
}

func allocTester(t *testing.T) unsafe.Pointer {
	size := rand.Intn(szpage * (1 + rand.Intn(8)))
	p := mem.Alloc(uint(size))
	b1 := *(*[]byte)(unsafe.Pointer(&slice{p, size, size}))
	b2 := make([]byte, size)
	rand.Read(b2)
	if n := copy(b1, b2); n != size {
		t.Errorf("incorrect number of copied. want: %v, got: %v", size, n)
	}
	if bytes.Compare(b1, b2) != 0 {
		t.Errorf("original does not match allocated slice. want: %x, got: %x", b2, b1)
	}
	return p
}

func TestConsecutive(t *testing.T) {
	var allocs []unsafe.Pointer
	for i := 0; i < 50; i++ {
		allocs = append(allocs, allocTester(t))
	}
	for _, p := range allocs {
		p := p
		if err := errpanic(func() { mem.Free(p) }); err != nil {
			t.Errorf("mem.Free(%p) panics: %v", p, err)
		}
	}
}

func TestMixed(t *testing.T) {
	var allocs []unsafe.Pointer
	for i := 0; i < 30; i++ {
		allocs = append(allocs, allocTester(t))
	}
	rand.Shuffle(len(allocs), func(i, j int) {
		allocs[i], allocs[j] = allocs[j], allocs[i]
	})
	for i := 29; i >= 10; i-- {
		p := allocs[i]
		if err := errpanic(func() { mem.Free(p) }); err != nil {
			t.Errorf("mem.Free(%p) panics: %v", p, err)
		}
		allocs = allocs[:i]
	}
	for i := 0; i < 20; i++ {
		allocs = append(allocs, allocTester(t))
	}
	for _, p := range allocs {
		p := p
		if err := errpanic(func() { mem.Free(p) }); err != nil {
			t.Errorf("mem.Free(%p) panics: %v", p, err)
		}
	}
}
