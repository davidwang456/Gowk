package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
)

func checkProcess(name string) bool {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("tasklist")
	case "linux", "darwin":
		cmd = exec.Command("ps", "-A")
	default:
		fmt.Printf("不支持的操作系统: %s\n", runtime.GOOS)
		return false
	}

	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("执行命令出错: %v\n", err)
		return false
	}

	processes := strings.Split(string(output), "\n")
	searchName := strings.ToLower(name)

	for _, process := range processes {
		processLower := strings.ToLower(process)
		// Windows 和 Linux 的进程名位置不同，需要分别处理
		switch runtime.GOOS {
		case "windows":
			// Windows 的 tasklist 输出格式：
			// "进程名.exe                          PID  ..."
			fields := strings.Fields(processLower)
			if len(fields) > 0 {
				procName := fields[0]
				// 移除可能的 .exe 后缀
				//procName = strings.TrimSuffix(procName, ".exe")
				if strings.Contains(procName, searchName) {
					return true
				}
			}
		case "linux", "darwin":
			// Linux/Mac 的 ps -A 输出格式：
			// "PID TTY          TIME CMD"
			fields := strings.Fields(processLower)
			if len(fields) >= 4 {
				procName := fields[3]
				if strings.Contains(procName, searchName) {
					return true
				}
			}
		}
	}
	return false
}

func processHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "只支持 GET 请求", http.StatusMethodNotAllowed)
		return
	}

	processName := r.URL.Query().Get("name")
	if processName == "" {
		http.Error(w, "请提供进程名称参数 'name'", http.StatusBadRequest)
		return
	}

	exists := checkProcess(processName)

	// 添加系统信息到响应中
	osInfo := fmt.Sprintf("当前操作系统: %s\n", runtime.GOOS)
	if exists {
		fmt.Fprintf(w, "%s进程 '%s' 正在运行\n", osInfo, processName)
	} else {
		fmt.Fprintf(w, "%s进程 '%s' 未运行\n", osInfo, processName)
	}
}

func main() {
	http.HandleFunc("/check", processHandler)

	fmt.Printf("服务器启动在 :9080 端口...(操作系统: %s)\n", runtime.GOOS)
	if err := http.ListenAndServe(":9080", nil); err != nil {
		fmt.Printf("服务器启动失败: %v\n", err)
	}
}
