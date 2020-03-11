// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris

package mem

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

func mmap(size uint) (unsafe.Pointer, error) {
	b, err := unix.Mmap(-1, 0, int(size), unix.PROT_READ|unix.PROT_WRITE, unix.MAP_PRIVATE|unix.MAP_ANONYMOUS)
	if err != nil {
		return nil, err
	}
	return unsafe.Pointer(&b[0]), nil
}

func munmap(p unsafe.Pointer) error {
	size := int(((*header)(p)).size + szheader)
	var sl = struct {
		addr unsafe.Pointer
		len  int
		cap  int
	}{p, size, size}
	b := *(*[]byte)(unsafe.Pointer(&sl))
	return unix.Munmap(b)
}
