package kid

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// GetMustMachineCode 获取设备机器码， 4位
// 有限使用IP后2位， 不存在使用Mac后2位， 不存在使用“0000”
func GetMustMachineCode() string {
	if ip := GetMustOneLocalIP(); ip != "" {
		ips := strings.Split(ip, ".")
		ip2, _ := strconv.Atoi(ips[2])
		ip3, _ := strconv.Atoi(ips[3])
		str := fmt.Sprintf("%02X%02X", ip2, ip3)
		return str
	} else if mac := GetMustOneLocalMac(); mac != "" {
		macs := strings.Split(mac, ":")
		str := fmt.Sprintf("%s%s", macs[len(macs)-2], macs[len(macs)-1])
		return strings.ToUpper(str)
	}
	return "0000"
}

// GetMustMachineCode2 获取设备机器码 4位
func GetMustMachineCode2() string {
	if mac := GetMustOneLocalMac(); mac != "" {
		macs := strings.Split(mac, ":")
		str := fmt.Sprintf("%s%s", macs[len(macs)-2], macs[len(macs)-1])
		return strings.ToUpper(str)
	}
	return "0000"
}

// GetMustOneLocalMac ...
func GetMustOneLocalMac() string {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("fail to get net interfaces: %v", err)
		return ""
	}

	for _, netInterface := range netInterfaces {
		macAddr := netInterface.HardwareAddr.String()
		if len(macAddr) == 0 {
			continue
		}
		return macAddr
	}
	return ""
}

// GetMustOneLocalIP ...
func GetMustOneLocalIP() string {
	ipInterfaces, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Printf("fail to get net interface addrs: %v", err)
		return ""
	}

	for _, address := range ipInterfaces {
		ipNet, ok := address.(*net.IPNet)
		if ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	return ""
}

// GetLocalMac ...
func GetLocalMac() (macAddrs []string) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("fail to get net interfaces: %v", err)
		return
	}

	for _, netInterface := range netInterfaces {
		macAddr := netInterface.HardwareAddr.String()
		if len(macAddr) == 0 {
			continue
		}
		macAddrs = append(macAddrs, macAddr)
	}
	return
}

// GetLocalIP ...
func GetLocalIP() (ips []string) {
	ipInterfaces, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Printf("fail to get net interface addrs: %v", err)
		return
	}

	for _, address := range ipInterfaces {
		ipNet, ok := address.(*net.IPNet)
		if ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP.String())
			}
		}
	}
	return ips
}
