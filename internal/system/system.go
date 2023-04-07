package system

import (
	"log"
	"net"
	"strings"
)

func GetIp() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func GetLocalIPs() []string {
	interfaces, _ := net.Interfaces()
	ips := []string{}
	for _, i := range interfaces {
		addresses, _ := i.Addrs()
		for _, address := range addresses {
			if strings.Split(address.String(), "/")[0] == GetIp() {
				ips = GenerateSubnetIPs(address.String())
			}
		}
	}
	return ips
}

func GenerateSubnetIPs(cidr string) []string {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return []string{}
	}
	ips := []string{}
	ip = ip.Mask(ipnet.Mask)
	for {
		lastIdx := len(ip) - 1
		ip[lastIdx]++
		ips = append(ips, ip.String())
		if ip[lastIdx] == 254 {
			break
		}
	}
	return ips
}
