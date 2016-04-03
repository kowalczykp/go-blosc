package blosc

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"io/ioutil"
	"testing"
	"unsafe"
)

func getTestInput() ([]byte, error) {
	input := new(bytes.Buffer)
	data := make([]float32, 100*100*100)
	// Generate hard to compress data
	for c := range data {
		data[c] = float32(c)
	}
	for _, c := range data {
		err := binary.Write(input, binary.LittleEndian, c)
		if err != nil {
			return nil, err
		}
	}
	return input.Bytes(), nil
}

func BenchmarkBloscSingleThread(b *testing.B) {
	var inputType float32
	output := ioutil.Discard
	input, err := getTestInput()
	if err != nil {
		b.Error(err)
		return
	}
	var size int
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		bw := NewWriter(output, unsafe.Sizeof(inputType))

		size, err = bw.Write(input)
		if err != nil {
			b.Error(err)
			return
		}
		bw.Close()
	}
	b.Logf("Result size: %d", size)
}

func BenchmarkBloscMultithread(b *testing.B) {
	var inputType float32
	output := ioutil.Discard
	input, err := getTestInput()
	if err != nil {
		b.Error(err)
		return
	}
	var size int
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		bw := NewWriter(output, unsafe.Sizeof(inputType))
		bw.SetMultithreading(4)
		size, err = bw.Write(input)
		if err != nil {
			b.Error(err)
			return
		}
		bw.Close()
	}
	b.Logf("Result size: %d", size)
}

func BenchmarkBloscSingleThreadZlib(b *testing.B) {
	var inputType float32
	output := ioutil.Discard
	input, err := getTestInput()
	if err != nil {
		b.Error(err)
		return
	}
	var size int
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		bw := NewWriter(output, unsafe.Sizeof(inputType))
		bw.SetCompressor("zlib")
		size, err = bw.Write(input)
		if err != nil {
			b.Error(err)
			return
		}
		bw.Close()
	}
	b.Logf("Result size: %d", size)
}

func BenchmarkBloscSingleThreadSnappy(b *testing.B) {
	var inputType float32
	output := ioutil.Discard
	input, err := getTestInput()
	if err != nil {
		b.Error(err)
		return
	}
	var size int
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		bw := NewWriter(output, unsafe.Sizeof(inputType))
		bw.SetCompressor("snappy")
		size, err = bw.Write(input)
		if err != nil {
			b.Error(err)
			return
		}
		bw.Close()
	}
	b.Logf("Result size: %d", size)
}

func BenchmarkGzip(b *testing.B) {
	output := ioutil.Discard
	input, err := getTestInput()
	if err != nil {
		b.Error(err)
		return
	}
	var size int

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		w := gzip.NewWriter(output)
		size, err = w.Write(input)
		if err != nil {
			b.Error(err)
			return
		}
		w.Close()
	}
	b.Logf("Result size: %d", size)

}
