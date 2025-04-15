package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	numLines = 5_000_000
	line     = "This is a line of text.\n"
)

func writeBuffered(filename string, bufSize int) time.Duration {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	writer := bufio.NewWriterSize(f, bufSize)

	start := time.Now()
	for i := 0; i < numLines; i++ {
		_, err := writer.WriteString(line)
		if err != nil {
			panic(err)
		}
	}
	writer.Flush()
	return time.Since(start)
}
func readBuffered(filename string, bufSize int) time.Duration {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	reader := bufio.NewReaderSize(f, bufSize)
	start := time.Now()

	for {
		_, err := reader.ReadString('\n')
		if err != nil {
			break
		}
	}
	return time.Since(start)
}
func main() {
	bufferSizes := []int{16, 256, 4096, 8192, 16384, 32768, 65536, 131072, 262144, 524288, 1048576} // 4KB to 1MB
	fmt.Printf("Running Buffer Size Tests: NumLines=%v", numLines)
	csvFile, err := os.Create("benchmark_results.csv")
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	writer.Write([]string{"BufferSize", "WriteTimeSeconds"}) // CSV header

	fmt.Println(" ----- Write Test ----- ")

	for _, size := range bufferSizes {
		duration := writeBuffered("temp.txt", size)
		fmt.Printf("Buffer size %d bytes: %v\n", size, duration)
		writer.Write([]string{strconv.Itoa(size), fmt.Sprintf("%.6f", duration.Seconds())})
	}
	fmt.Println(" ----- Read Test ----- ")
	for _, size := range bufferSizes {
		duration := readBuffered("temp.txt", size)
		fmt.Printf("Buffer size %d bytes: %v\n", size, duration)
		writer.Write([]string{strconv.Itoa(size), fmt.Sprintf("%.6f", duration.Seconds())})
	}
}
