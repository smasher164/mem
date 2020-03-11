// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris

package mem

import (
	"reflect"
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
	var b []byte
	sl := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sl.Data = uintptr(p)
	sl.Len = size
	sl.Cap = size
	return unix.Munmap(b)
}
