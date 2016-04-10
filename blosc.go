package blosc

/*
#cgo LDFLAGS: -lblosc

#include <blosc.h>
*/
import "C"
import (
	"errors"
	"unsafe"
)

const (
	MaxOverhead   int = C.BLOSC_MAX_OVERHEAD
	MaxBufferSize int = C.BLOSC_MAX_BUFFERSIZE
)

type shuffle int

const (
	NoShuffle  shuffle = C.BLOSC_NOSHUFFLE
	Shuffle    shuffle = C.BLOSC_SHUFFLE
	BitShuffle shuffle = C.BLOSC_BITSHUFFLE
)

type compressor string

const (
	Blosclz compressor = "blosclz"
	Lz4     compressor = "lz4"
	Lz4hc   compressor = "lz4hc"
	Snappy  compressor = "snappy"
	Zlib    compressor = "zlib"
)

type SizeInfo struct {
	Original   int
	Compressed int
	Block      int
}

func GetSizeInfo(src []byte) (*SizeInfo, error) {
	if src == nil {
		return nil, errors.New("blosc: empty src")
	}
	var oSize, cSize, bSize C.size_t
	C.blosc_cbuffer_sizes(unsafe.Pointer(&src[0]), &oSize, &cSize, &bSize)
	return &SizeInfo{Original: int(oSize), Compressed: int(cSize), Block: int(bSize)}, nil
}
