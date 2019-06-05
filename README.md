# mem

Package mem implements a memory allocator and deallocator.
It currently uses mmap on unix and VirtualAlloc on windows
to request pages of memory from the operating system, and
munmap and VirtualFree to release pages of memory to the
operating system. The allocator uses a first-fit algorithm
on a singly-linked free list of blocks. Blocks are divided
into sets called arenas, which correspond to the chunk of
memory mapped from the operating system. When all of the
blocks in a set are freed, the arena is unmapped.

```
// Alloc allocates size bytes of memory, and returns a pointer to it.
// It is goroutine-safe and attempts to preserve the semantics of
// POSIX libc's malloc. However, Alloc panics if an error occurs when
// requesting more memory from the operating system.
func Alloc(size uint) unsafe.Pointer

// Free deallocates the memory pointed to by p. It is goroutine-safe
// and attempts to preserve the semantics of POSIX libc's free.
// However, Free panics if an error occurs when releasing memory to
// the operating system.
func Free(p unsafe.Pointer)
```
### Why?

I am working on a rudimentary language interpreter for which I want
to implement my own garbage collection algorithm in Go. Also, this
is a good opportunity for me to brush up on memory allocation
techniques.

### Should I use this?

If you want to? ¯\\\_(ツ)\_/¯ It's passed all the tests I've run so far.
But don't hold me liable if your multi-million dollar production
workload comes crashing down! Instead, file an issue. :)

### What can I do to help?

Why thank you for offering! I am really interested in exploring
more efficient memory allocation techniques, like using a
best-fit algorithm, using a doubly-linked list for constant-time
cleanup, segregated lists to divide allocations of different
size classes, and buddy allocation. I would also like to improve
test quality by testing edge cases and platform-specific eccentricities.