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
	C.blosc_init()
	return NewLevelWriter(w, 5, Shuffle, size)
}

func NewLevelWriter(w io.Writer, level int, s shuffle, size uintptr) *Writer {
	return &Writer{w: w, level: level, size: size, s: s}
}

type Writer struct {
	w     io.Writer
	level int
	size  uintptr
	s     shuffle
}

func (w *Writer) Close() error {
	C.blosc_destroy()
	return nil
}

func (w *Writer) Flush() error {
	return nil
}

func (w *Writer) Reset(wr io.Writer) {
	return
}

func (w *Writer) Write(p []byte) (n int, err error) {
	dest := make([]byte, len(p))
	csize := C.blosc_compress(C.int(w.level), C.int(w.s), C.size_t(w.size), C.size_t(len(dest)), unsafe.Pointer(&p[0]), unsafe.Pointer(&dest[0]), C.size_t(len(dest)))
	if csize == 0 {
		return 0, errors.New("uncompressible buffer")
	} else if csize < 0 {
		return int(csize), errors.New(fmt.Sprintf("compression error %d", csize))
	}
	return w.w.Write(dest[:int(csize)])
}
