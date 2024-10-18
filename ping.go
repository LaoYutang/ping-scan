package main

import (
	"net"
	"os/exec"
	"runtime"
)

// Ping 测试指定 IP 是否连通
func Ping(ipAddr string) bool {
	var cmd *exec.Cmd

	// 根据操作系统选择适合的 ping 命令
	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-n", "1", "-w", "200", ipAddr)
	} else {
		cmd = exec.Command("ping", "-c", "1", "-W", "0.2", ipAddr)
	}

	err := cmd.Run()
	return err == nil
}

// IPRange 生成指定网段中的所有 IP 地址
func IPRange(cidr string) ([]string, error) {
	var ips []string
	// 解析 CIDR
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}

	return ips, nil
}

// inc 用于递增 IP 地址
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
