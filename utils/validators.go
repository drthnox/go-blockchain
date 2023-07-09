package utils

import (
	"fmt"
	"net"
)

func CheckIPAddress(ip string) bool {
	if net.ParseIP(ip) == nil {
		fmt.Printf("IP Address: %s - Invalid\n", ip)
		return false
	} else {
		fmt.Printf("IP Address: %s - Valid\n", ip)
		return true
	}
}
