package remarkablepage

/*
#cgo CFLAGS: -I. -ffast-math
#cgo LDFLAGS: -L. -limage_processing -lm -lpthread -fopenmp -O3 -mfpu=neon -march=armv7-a
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

const maxSize = 1 << 16 // 2¹⁶

func HandleNewFile(directory, filename string) LineList {
	dir := C.CString(directory)
	file := C.CString(filename)
	defer C.free(unsafe.Pointer(dir))
	defer C.free(unsafe.Pointer(file))
	ll := C.handle_new_file(dir, file)

	size := int(ll.size)
	if size > maxSize {
		size = maxSize
	}
	lines := (*[maxSize]float32)(unsafe.Pointer(ll.lines))[: size*4 : size*4]

	return LineList{Lines: lines, Size: size}
}
