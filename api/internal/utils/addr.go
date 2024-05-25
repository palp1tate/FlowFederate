package utils

import (
	"net"
	"os"
	"runtime"
)

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func GetIPAddress() string {
	switch runtime.GOOS {
	case "windows":
		// 在Windows上获取本机名
		hostname, err := os.Hostname()
		if err != nil {
			return "127.0.0.1" // 如果无法获取hostname，返回本地地址
		}
		// 查找hostname对应的IP地址
		addrs, err := net.LookupIP(hostname)
		if err != nil || len(addrs) == 0 {
			return "127.0.0.1"
		}
		// 返回第一个IPv4地址
		for _, addr := range addrs {
			if ipv4 := addr.To4(); ipv4 != nil {
				return ipv4.String()
			}
		}
		return "127.0.0.1"
	case "linux", "darwin":
		// 在Linux或macOS上获取与外部通信的接口的IP地址
		conn, err := net.Dial("udp", "8.8.8.8:80")
		if err != nil {
			return "127.0.0.1"
		}
		defer conn.Close()
		localAddr := conn.LocalAddr().(*net.UDPAddr)
		return localAddr.IP.String()
	default:
		return "Unsupported OS"
	}
}
