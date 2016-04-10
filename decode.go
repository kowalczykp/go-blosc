package blosc

/*
#cgo LDFLAGS: -lblosc

#include <blosc.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// TODO: Reader implementation - needs access to internal functions in blosc library

type Decoder struct {
	threads C.int
}

func NewDecoder() *Decoder {
	return NewThreadedDecoder(1)
}

func NewThreadedDecoder(threads int) *Decoder {
	return &Decoder{threads: C.int(threads)}
}

func (d *Decoder) Decode(src, dst []byte) ([]byte, error) {
	meta, err := GetSizeInfo(src)
	if err != nil {
		return nil, err
	}
	if len(dst) < meta.Original {
		dst = make([]byte, meta.Original)
	}
	csize := C.blosc_decompress_ctx(unsafe.Pointer(&src[0]), unsafe.Pointer(&dst[0]), C.size_t(len(dst)), d.threads)
	if csize < 1 {
		return nil, fmt.Errorf("blosc: decompression error %d", csize)
	}
	return dst[:int(csize)], nil
}
