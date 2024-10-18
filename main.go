package main

import (
	"bytes"
	"fmt"
	"net"
	"runtime"
	"sort"
	"sync"
	"time"
)

func main() {
	fmt.Println("请输入要 ping 的网段(xxx:xxx:xxx:xxx/xx)：")
	var cidr string
	fmt.Scanln(&cidr)

	ips, err := IPRange(cidr)
	if err != nil {
		fmt.Println("Error generating IP range:", err)
		return
	}

	// 获取 CPU 核心数，限制协程数量为 CPU 核心数的8倍
	sem := make(chan struct{}, runtime.NumCPU()*8)

	// 创建 WaitGroup 以等待所有 goroutine 完成
	var wg sync.WaitGroup
	wg.Add(len(ips))

	// 创建一个通道用于输出成功的结果
	resultCh := make(chan string, len(ips))
	progressCh := make(chan int, len(ips))

	start := time.Now()

	// 协程监控进度
	go func() {
		total := len(ips)
		processed := 0
		for p := range progressCh {
			processed += p
			fmt.Printf("\r当前进度: %d/%d", processed, total) // 打印进度
		}
	}()

	// 多协程 Ping 测试每个 IP
	for _, ip := range ips {
		sem <- struct{}{} // 占用一个协程
		go func(ip string) {
			defer wg.Done()
			defer func() { <-sem }() // 释放一个协程

			if Ping(ip) {
				resultCh <- ip // 仅将成功的 IP 传入通道
			}
			progressCh <- 1 // 每次 ping 完成后更新进度
		}(ip)
	}

	// 等待所有 goroutine 完成
	go func() {
		wg.Wait()
		close(resultCh)
		close(progressCh)
	}()

	// 收集成功 ping 的 IP 地址
	var successfulIPs []net.IP
	for ip := range resultCh {
		successfulIPs = append(successfulIPs, net.ParseIP(ip))
	}

	// 使用自定义排序
	sort.Slice(successfulIPs, func(i, j int) bool {
		return bytes.Compare(successfulIPs[i], successfulIPs[j]) < 0
	})

	// 输出 ping 成功的 IP 地址
	fmt.Println("\nPing 成功的 IP 地址如下：")
	for _, ip := range successfulIPs {
		fmt.Println(ip)
	}

	fmt.Printf("Ping测试总计用时 %s\n", time.Since(start))

	// 等待任意键退出
	fmt.Println("按任意键退出...")
	fmt.Scanln()
}
