package Core

import (
	"net"
	"strconv"
)

type NetInfo struct {
	//每次最好尝试找一个端口给远程服务
	Port string
}

func (n NetInfo) GetNetIPs() ([]string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	var result []string
	result = append(result, "127.0.0.1")
	for _, addr := range addrs {
		// 检查IP地址的类型，只获取IPv4地址
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				// fmt.Println("IP地址:", ipnet.IP.String())
				result = append(result, ipnet.IP.String())
			}
		}
	}
	return result, nil
}

func NewNetInfo() NetInfo {
	// port := 50520
	// count := 100
	// for count > 0 {
	// 	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	// 	if err != nil {
	// 		fmt.Printf("Remote Port %d is already in use\n", port)
	// 		port++
	// 	} else {
	// 		listener.Close()
	// 		fmt.Printf("Remote Port %d is available\n", port)
	// 		break
	// 	}
	// 	count--
	// }
	return NetInfo{Port: strconv.Itoa(50520)}
}
