package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unsafe"

	"github.com/kowalczykp/go-blosc"
)

func main() {
	input := new(bytes.Buffer)
	output := new(bytes.Buffer)

	data := make([]float32, 100*100*100)

	bw := blosc.NewWriter(output, unsafe.Sizeof(data[0]))
	defer bw.Close()

	// Generate hard to compress data
	for c := range data {
		data[c] = float32(c)
	}

	for _, c := range data {
		err := binary.Write(input, binary.LittleEndian, c)
		if err != nil {
			fmt.Println("binary.Write failed:", err)
			return
		}
	}

	written, err := bw.Write(input.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Compression: %d -> %d (%.1fx)\n", len(data), written, float32(len(data))/float32(written))
}
