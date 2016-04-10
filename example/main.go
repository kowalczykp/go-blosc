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
	data := make([]float32, 100*100*100)

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

	originalSize := len(input.Bytes())

	// Compression
	enc := blosc.NewEncoder()
	compressed, err := enc.Encode(unsafe.Sizeof(data[0]), input.Bytes(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	compressedSize := len(compressed)
	ratio := float32(originalSize) / float32(compressedSize)
	fmt.Printf("Original size: %d after compression: %d Ratio: %f\n", originalSize, compressedSize, ratio)

	// Decompression
	dec := blosc.NewDecoder()
	decompressed, err := dec.Decode(compressed, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Original==Decompressed:", bytes.Compare(decompressed, input.Bytes()) == 0)

}
