package mem

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

func mmap(size uint) (unsafe.Pointer, error) {
	p, err := windows.VirtualAlloc(0, uintptr(size), windows.MEM_COMMIT|windows.MEM_RESERVE, windows.PAGE_READWRITE)
	if err != nil {
		return nil, err
	}
	return unsafe.Pointer(p), nil
}

func munmap(p unsafe.Pointer) error {
	return windows.VirtualFree(uinptr(p), 0, windows.MEM_RELEASE)
}
