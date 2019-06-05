// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris

package mem

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

type slice struct {
	addr unsafe.Pointer
	len  int
	cap  int
}

func mmap(size uint) (unsafe.Pointer, error) {
	b, err := unix.Mmap(-1, 0, int(size), unix.PROT_READ|unix.PROT_WRITE, unix.MAP_PRIVATE|unix.MAP_ANONYMOUS)
	if err != nil {
		return nil, err
	}
	sl := (*slice)(unsafe.Pointer(&b))
	return sl.addr, nil
}

func munmap(p unsafe.Pointer) error {
	size := int(((*header)(p)).size + szheader)
	b := *(*[]byte)(unsafe.Pointer(&slice{p, size, size}))
	return unix.Munmap(b)
}
