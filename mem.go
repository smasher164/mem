// Package mem implements a memory allocator and deallocator.
// It currently uses mmap on unix and VirtualAlloc on windows
// to request pages of memory from the operating system, and
// munmap and VirtualFree to release pages of memory to the
// operating system. The allocator uses a first-fit algorithm
// on a singly-linked free list of blocks. Blocks are divided
// into sets called arenas, which correspond to the chunk of
// memory mapped from the operating system. When all of the
// blocks in a set are freed, the arena is unmapped.
package mem

import (
	"os"
	"sync"
	"unsafe"
)

var (
	szpage   = uint(os.Getpagesize())
	szheader = uint(unsafe.Sizeof(header{}))
	freep    = new(header) // head of linked list
	m        sync.Mutex
)

type header struct {
	size      uint
	allocated bool
	next      *header
	arena     unsafe.Pointer // arena of which this block is a member
}

// return n rounded up to a multiple of k.
func roundMultiple(n, k uint) uint {
	if mod := n % k; mod != 0 {
		return n + (k - mod)
	}
	return n
}

// Alloc allocates size bytes of memory, and returns a pointer to it.
// It is goroutine-safe and attempts to preserve the semantics of
// POSIX libc's malloc. However, Alloc panics if an error occurs when
// requesting more memory from the operating system.
func Alloc(size uint) unsafe.Pointer {
	if size == 0 {
		return nil
	}
	m.Lock()
	defer m.Unlock()
	// Iterate through linked list of headers.
	p := freep
	for {
		// block is free
		if !p.allocated {
			// block is large enough (first-fit)
			if p.size >= size {
				ret := uintptr(unsafe.Pointer(p)) + uintptr(szheader)
				// split block if enough space for header
				if gap := p.size - size; gap >= szheader {
					h := (*header)(unsafe.Pointer(ret + uintptr(size)))
					h.size = gap - szheader
					h.allocated = false
					h.next = p.next
					h.arena = p.arena
					p.next = h
					p.size = size
				}
				p.allocated = true
				return unsafe.Pointer(ret)
			}
		}
		// allocate memory
		if p.next == nil {
			// Allocated space must be enough to hold header and
			// size bytes. Simplify alignment by rounding up to
			// the next multiple of the header size.
			szalign := roundMultiple(szheader+size, szheader)
			// round bytes to the next multiple of the page size
			szalloc := roundMultiple(szalign, szpage)
			pblock, err := mmap(szalloc)
			if err != nil {
				panic(err)
			}
			h := (*header)(pblock)
			h.size = szalloc - szheader
			h.allocated = false
			h.next = nil
			h.arena = pblock
			p.next = h
		}
		p = p.next
	}
}

// Free deallocates the memory pointed to by p. It is goroutine-safe
// and attempts to preserve the semantics of POSIX libc's free.
// However, Free panics if an error occurs when releasing memory to
// the operating system.
func Free(p unsafe.Pointer) {
	if p == nil {
		return
	}
	m.Lock()
	defer m.Unlock()
	h := (*header)(unsafe.Pointer(uintptr(p) - uintptr(szheader)))
	arcurr := h.arena
	h.allocated = false
	// Coalesce adjacent free blocks from the same arena.
	if next := h.next; next != nil && next.arena == h.arena && !next.allocated {
		h.size += next.size + szheader
		h.next = next.next
	}
	freeArena := true
	var first, prev *header
	for it := freep; ; it = it.next {
		if it.arena == arcurr {
			if it.allocated {
				freeArena = false
			} else if next := it.next; next != nil && next.arena == arcurr && !next.allocated {
				it.size += next.size + szheader
				it.next = next.next
			}
			if first == nil {
				first = it
			}
		}
		if it.next == nil {
			break
		}
		if it.arena != arcurr && it.next.arena == arcurr {
			prev = it
		}
	}
	if !freeArena {
		return
	}
	// If there are no allocated blocks in arena, munmap it.
	prev.next = first.next
	if err := munmap(unsafe.Pointer(first)); err != nil {
		panic(err)
	}
}
