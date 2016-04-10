package blosc

/*
#cgo LDFLAGS: -lblosc

#include <blosc.h>
*/
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

// TODO: Writer implementation - needs access to internal functions in blosc library

type Encoder struct {
	level      C.int
	shuffle    C.int
	threads    C.int
	compressor string
}

func NewEncoder() *Encoder {
	return NewAdvancedEncoder(5, Shuffle, 1, Blosclz)
}

func NewAdvancedEncoder(level int, s shuffle, threads int, comp compressor) *Encoder {
	return &Encoder{level: C.int(level), shuffle: C.int(s), threads: C.int(threads), compressor: string(comp)}
}

func (e *Encoder) Encode(tSize uintptr, src, dst []byte) ([]byte, error) {
	out := dst
	if len(out) < len(src)+MaxOverhead-1 {
		out = make([]byte, len(src)+MaxOverhead)
	}
	// TODO: replace with a maping to C-allocated strings, to avoid allocations and frees here
	cCompressor := C.CString(e.compressor)
	defer C.free(unsafe.Pointer(cCompressor))

	csize := C.blosc_compress_ctx(e.level, e.shuffle, C.size_t(tSize),
		C.size_t(len(src)), unsafe.Pointer(&src[0]), unsafe.Pointer(&out[0]), C.size_t(len(out)),
		cCompressor, 0, e.threads)
	if csize == 0 {
		return nil, errors.New("blosc: uncompressible buffer")
	} else if csize < 0 {
		return nil, fmt.Errorf("blosc: compression error %d", csize)
	}
	return out[:int(csize)], nil
}
