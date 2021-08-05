package common

import (
	"log"
	"net"
	"os"
)

// 获取内网IP
func GetInternalIp() string {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		log.Println("net.Interfaces failed, err:", err.Error())
		return "127.0.0.1"
	}
	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()
			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String()
					}
				}
			}
		}
	}
	return "127.0.0.1"
}

func IsDirExists(fileAddr string)bool{
	s,err:=os.Stat(fileAddr)
	if err!=nil{
		log.Println(err)
		return false
	}
	return s.IsDir()
}
