package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/pipeline"
)

func main() {
	p := CreatePipelne("larg.in", 400, 4)
	WriteToFile(p, "small.out")
	PrintFile("small.out")
}

// CreatePipelne ...
// chunkCount 分块的数目
func CreatePipelne(filename string, fileSize, chunkCount int) <-chan int {
	chunkSize := fileSize / chunkCount

	sortResults := []<-chan int{} // 收集文件中的分块
	for i := 0; i < chunkCount; i++ {
		file, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		// offset int64: i*chunkSize
		// 0: the origin of the file
		file.Seek(int64(i*chunkSize), 0)

		source := pipeline.ReadSource(bufio.NewReader(file), chunkSize)

		sortResults = append(sortResults, pipeline.InMemSort(source))
	}
	// 将分块进行归并
	return pipeline.MergeN(sortResults...)
}

// WriteToFile ...
func WriteToFile(p <-chan int, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	pipeline.WriteSink(writer, p)
}

// PrintFile ...
func PrintFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// 通过它，我们可以从底层的 io.Reader 中更大批量的读取数据。这会使读取操作变少。如果数据读取时的块数量是固定合适的，
	// 底层媒体设备将会有更好的表现，也因此会提高程序的性能：
	// io.Reader --> buffer --> consumer
	p := pipeline.ReadSource(bufio.NewReader(file), -1)
	count := 0
	for i := range p {
		fmt.Println(i)
		count++
		if count >= 20 {
			break
		}
	}
}
