package blosc

/*
#cgo LDFLAGS: -lblosc

#include <blosc.h>

*/
import "C"

import (
	"errors"
	"fmt"
	"io"
	"unsafe"
)

type shuffle int

const (
	NoShuffle  shuffle = 0
	Shuffle    shuffle = 1
	BitShuffle shuffle = 2
)

func NewWriter(w io.Writer, size uintptr) *Writer {
	return NewLevelWriter(w, 5, Shuffle, size)
}

func NewLevelWriter(w io.Writer, level int, s shuffle, size uintptr) *Writer {
	C.blosc_init()
	return &Writer{w: w, level: level, size: size, s: s}
}

type Writer struct {
	w     io.Writer
	level int
	size  uintptr
	s     shuffle
}

func (b *Writer) SetMultithreading(n int) {
	C.blosc_set_nthreads(C.int(n))
}

func (b *Writer) SetCompressor(compressor string) {
	cCompressor := C.CString(compressor)
	C.blosc_set_compressor(cCompressor)
	C.free(unsafe.Pointer(cCompressor))
}

func (b *Writer) Close() error {
	C.blosc_destroy()
	return nil
}

func (b *Writer) Flush() error {
	return nil
}

func (b *Writer) Reset(wr io.Writer) {
	return
}

func (b *Writer) Write(p []byte) (n int, err error) {
	dest := make([]byte, len(p))
	csize := C.blosc_compress(C.int(b.level), C.int(b.s), C.size_t(b.size), C.size_t(len(dest)), unsafe.Pointer(&p[0]), unsafe.Pointer(&dest[0]), C.size_t(len(dest)))
	if csize == 0 {
		return 0, errors.New("uncompressible buffer")
	} else if csize < 0 {
		return int(csize), errors.New(fmt.Sprintf("compression error %d", csize))
	}
	return b.w.Write(dest[:int(csize)])
}
