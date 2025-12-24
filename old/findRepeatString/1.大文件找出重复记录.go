package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

func main() {
	file1 := "file1.txt"
	file2 := "file2.txt"
	outputFile := "duplicates.txt"

	// 创建临时目录
	tmpDir := "./tmp_chunks"
	os.RemoveAll(tmpDir)
	os.MkdirAll(tmpDir, 0755)

	// 阶段1：哈希分片处理两个文件
	chunkFiles1 := processFile(file1, tmpDir, "file1")
	chunkFiles2 := processFile(file2, tmpDir, "file2")

	// 阶段2：找出重复记录
	findDuplicates(chunkFiles1, chunkFiles2, outputFile, tmpDir)

	// 清理临时文件
	os.RemoveAll(tmpDir)

	fmt.Println("完成！重复记录已保存到", outputFile)
}

// 处理单个文件：哈希分片
func processFile(filename, tmpDir, prefix string) []string {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 创建分片文件
	chunkFiles := make([]*os.File, 256)
	chunkWriters := make([]*bufio.Writer, 256)
	chunkFileNames := make([]string, 256)

	for i := 0; i < 256; i++ {
		chunkName := filepath.Join(tmpDir, fmt.Sprintf("%s_chunk_%03d.txt", prefix, i))
		chunkFiles[i], _ = os.Create(chunkName)
		chunkWriters[i] = bufio.NewWriter(chunkFiles[i])
		chunkFileNames[i] = chunkName
	}
	defer func() {
		for i := 0; i < 256; i++ {
			if chunkWriters[i] != nil {
				chunkWriters[i].Flush()
			}
			if chunkFiles[i] != nil {
				chunkFiles[i].Close()
			}
		}
	}()

	scanner := bufio.NewScanner(file)
	buffer := make([]byte, 0, 64*1024)
	scanner.Buffer(buffer, 10*1024*1024) // 支持长行

	var wg sync.WaitGroup
	ch := make(chan string, 10000)

	// 启动多个worker并行处理
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for line := range ch {
				hash := sha256.Sum256([]byte(line))
				chunkIndex := int(hash[0]) // 用第一个字节作为分片索引

				chunkWriters[chunkIndex].WriteString(line + "\n")
			}
		}()
	}

	// 读取文件并分发
	for scanner.Scan() {
		ch <- scanner.Text()
	}
	close(ch)
	wg.Wait()

	// 刷新所有writer
	for i := 0; i < 256; i++ {
		chunkWriters[i].Flush()
	}

	return chunkFileNames
}

// 找出重复记录
func findDuplicates(chunks1, chunks2 []string, outputFile, tmpDir string) {
	outFile, _ := os.Create(outputFile)
	defer outFile.Close()
	writer := bufio.NewWriter(outFile)
	defer writer.Flush()

	var wg sync.WaitGroup
	results := make(chan string, 1000)

	// 启动结果收集器
	go func() {
		for dup := range results {
			writer.WriteString(dup + "\n")
		}
	}()

	// 并行处理每个分片
	for i := 0; i < 256; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			// 读取第一个文件的分片到map
			set1 := make(map[string]bool)
			if f1, err := os.Open(chunks1[idx]); err == nil {
				scanner := bufio.NewScanner(f1)
				for scanner.Scan() {
					set1[scanner.Text()] = true
				}
				f1.Close()
			}

			// 检查第二个文件的分片
			if f2, err := os.Open(chunks2[idx]); err == nil {
				scanner := bufio.NewScanner(f2)
				for scanner.Scan() {
					line := scanner.Text()
					if set1[line] {
						results <- line
					}
				}
				f2.Close()
			}
		}(i)
	}

	wg.Wait()
	close(results)
}
