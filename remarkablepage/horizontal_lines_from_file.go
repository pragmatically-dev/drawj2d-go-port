package remarkablepage

/*
#cgo CFLAGS: -I. -ffast-math
#cgo LDFLAGS: -L.  -lm -lpthread -fopenmp -O3 -mfpu=neon -march=armv7-a
#include "image_processing.h"
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"
)

type LineList struct {
	Lines []float32
	Size  int
}

const maxSize = 1 << 28 // 2^(28)

func HandleNewFile(directory, filename string) LineList {
	dir := C.CString(directory)
	file := C.CString(filename)
	defer C.free(unsafe.Pointer(dir))
	defer C.free(unsafe.Pointer(file))

	// Ensure ll is properly allocated and not moved by GC
	ll := C.handle_new_file(dir, file)
	defer C.free(unsafe.Pointer(ll.lines))

	size := int(ll.size)
	if size > maxSize {
		size = maxSize
	}

	// Ensure ll.lines is properly pinned
	lines := make([]float32, size*4)
	copy(lines, (*[maxSize]float32)(unsafe.Pointer(ll.lines))[:size*4:size*4])

	return LineList{Lines: lines, Size: size}
}
