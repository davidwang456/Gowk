package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func dd() {
	// 配置参数
	totalLines := 1_000_000
	targetLines := 5
	fileName := "mylog.log"

	// 创建文件
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("创建文件错误: %v\n", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// 生成包含"File Name"的5个随机位置（避开开头和结尾的10000行）
	rand.Seed(time.Now().UnixNano())
	positions := make(map[int]bool)
	for len(positions) < targetLines {
		pos := rand.Intn(totalLines-20000) + 10000
		positions[pos] = true
	}

	// 写入日志
	for i := 0; i < totalLines; i++ {
		var line string
		if positions[i] {
			line = fmt.Sprintf("[%d] INFO: Processing File Name: document_%d.txt\n", i, rand.Intn(1000))
		} else {
			line = fmt.Sprintf("[%d] DEBUG: System status check completed. Status: OK\n", i)
		}
		if _, err := writer.WriteString(line); err != nil {
			fmt.Printf("写入文件错误: %v\n", err)
			return
		}

		// 每10000行刷新一次缓冲区
		if i%10000 == 0 {
			writer.Flush()
		}
	}

	fmt.Printf("日志文件生成完成！\n")
	fmt.Printf("总行数: %d\n", totalLines)
	fmt.Printf("包含'File Name'的行数: %d\n", targetLines)
	fmt.Println("包含关键词的行号：")
	for pos := range positions {
		fmt.Printf("第 %d 行\n", pos)
	}
}
