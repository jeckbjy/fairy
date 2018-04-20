package util

import (
	"errors"
	"net"
)

// GetIPv4 返回一个内网IPv4地址
func GetIPv4() (string, error) {
	as, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, a := range as {
		ipnet, ok := a.(*net.IPNet)
		if !ok || ipnet.IP.IsLoopback() {
			continue
		}

		ip := ipnet.IP.To4()
		if ip != nil {
			return ip.String(), nil
		}
	}

	return "", errors.New("no private ip address")
}
