package common

import (
	"bufio"
	"io"
	"net"
	"os"
)

// GetInternalIp 获取内网IP
func GetInternalIp() string {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		Log.Info("net.Interfaces failed, err:", err.Error())
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

// IsDirExists 目录是否存在
func IsDirExists(fileAddr string)bool{
	s,err:=os.Stat(fileAddr)
	if err!=nil{
		return false
	}
	return s.IsDir()
}

func IsFileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}


func ReadLog(fileAddr string, fromLineNum int) (string, int, bool) {
	logContent := ""
	lineNum := 0
	if !IsFileExist(fileAddr) {
		return "readLog fail, logFile not exists", 0, true
	}
	fd, err := os.Open(fileAddr)
	defer func(fd *os.File) {
		err := fd.Close()
		if err != nil {
			Log.Error(err)
		}
	}(fd)
	if err != nil {
		Log.Error("read error:", err)
	}
	buff := bufio.NewReader(fd)
	for {
		data, _, eof := buff.ReadLine()
		if eof == io.EOF {
			break
		}
		lineNum++
		if lineNum >= fromLineNum {
			logContent = logContent + (string(data)) + "\n"
		}
	}
	return logContent, lineNum, false
}
