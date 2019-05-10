package pipeline

import (
	"encoding/binary"
	"io"
	"math/rand"
	"sort"
)

//ArraySource send int to chan
func ArraySource(a ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, v := range a {
			out <- v
		}
	}()
	return out
}

// InMemSort ...
func InMemSort(a <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		s := []int{}
		for i := range a {
			s = append(s, i)
		}

		sort.Ints(s)

		for _, j := range s {
			out <- j
		}
	}()
	return out
}

// Merge ...
func Merge(a, b <-chan int) <-chan int {
	out := make(chan int)
	go func() {

		for {
			v1, ok1 := <-a
			v2, ok2 := <-b
			defer close(out)
			for ok1 || ok2 {
				if !ok2 || (ok1 && ok2 && v1 <= v2) {
					out <- v1
					v1, ok1 = <-a
				} else {
					out <- v2
					v2, ok2 = <-b
				}
			}
		}
	}()
	return out
}

// RandomSource ...
func RandomSource(count int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := 0; i < count; i++ {
			out <- rand.Int()
		}
	}()
	return out
}

// WriteSink ...
func WriteSink(writer io.Writer, in <-chan int) {
	for v := range in {
		buffer := make([]byte, 8)
		binary.BigEndian.PutUint64(buffer, uint64(v))
		writer.Write(buffer)
	}
}

// ReadSource ...
func ReadSource(reader io.Reader, chunkSize int) <-chan int {
	out := make(chan int)
	go func() {
		buffer := make([]byte, 8)
		bytesReader := 0
		defer close(out)
		for {
			n, err := reader.Read(buffer)
			bytesReader += n
			if n > 0 {
				i := int(binary.BigEndian.Uint64(buffer))
				out <- i
			}
			// 可能出错了 但是还读到了几个字节 不能扔掉
			if err != nil || (bytesReader > chunkSize && chunkSize != -1) {
				break
			}
		}
	}()
	return out
}

// MergeN ...
// 退出函数递归之后, 才逆顺序执行
func MergeN(inputs ...<-chan int) <-chan int {
	if len(inputs) == 1 {
		return inputs[0]
	}
	mid := len(inputs) / 2
	return Merge(MergeN(inputs[:mid]...), MergeN(inputs[mid:]...))
}
