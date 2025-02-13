package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

func main() {
	start := time.Now()

	// 配置参数
	inputFile := "mylog.log" // 修改为刚生成的日志文件
	outputFile := "output.csv"
	searchTerm := "File Name"
	numWorkers := runtime.NumCPU() // 使用CPU核心数作为工作协程数

	// 获取文件大小
	fileInfo, err := os.Stat(inputFile)
	if err != nil {
		fmt.Printf("获取文件信息错误: %v\n", err)
		return
	}

	// 创建输出CSV文件
	csvFile, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("创建CSV文件错误: %v\n", err)
		return
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// 创建通道用于传递匹配的行
	matches := make(chan string, 1000)
	var wg sync.WaitGroup

	// 计算每个工作协程处理的文件块大小
	chunkSize := fileInfo.Size() / int64(numWorkers)

	// 启动工作协程
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			// 为每个协程创建独立的文件句柄
			file, err := os.Open(inputFile)
			if err != nil {
				fmt.Printf("打开文件错误: %v\n", err)
				return
			}
			defer file.Close()

			// 计算该协程处理的文件范围
			startPos := int64(workerID) * chunkSize
			endPos := startPos + chunkSize
			if workerID == numWorkers-1 {
				endPos = fileInfo.Size()
			}

			// 设置文件读取位置
			if _, err := file.Seek(startPos, 0); err != nil {
				fmt.Printf("设置文件位置错误: %v\n", err)
				return
			}

			reader := bufio.NewReaderSize(file, 1024*1024) // 1MB 缓冲区

			// 如果不是第一个块，读取并丢弃第一个不完整的行
			if startPos > 0 {
				_, err := reader.ReadString('\n')
				if err != nil {
					fmt.Printf("读取错误: %v\n", err)
					return
				}
			}

			currentPos := startPos
			// 读取并处理数据
			for {
				line, err := reader.ReadString('\n')
				if err != nil {
					break
				}

				currentPos += int64(len(line))

				// 如果不是最后一个块，只读取到块的结尾
				if workerID < numWorkers-1 && currentPos > endPos {
					break
				}

				if strings.Contains(line, searchTerm) {
					matches <- strings.TrimSpace(line)
				}
			}
		}(i)
	}

	// 启动写入协程
	go func() {
		wg.Wait()
		close(matches)
	}()

	// 写入CSV
	lineCount := 0
	for match := range matches {
		writer.Write([]string{match})
		lineCount++
	}

	elapsed := time.Since(start)
	fmt.Printf("处理完成！共找到 %d 行匹配内容\n", lineCount)
	fmt.Printf("处理时间: %v\n", elapsed)
}
