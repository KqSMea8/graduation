package tools

import "net"

func GetCurrentIp() int32 {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
}
