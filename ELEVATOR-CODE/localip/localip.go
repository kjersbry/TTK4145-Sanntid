package localip

import (
	"net"
	"strings"
	"fmt"
)

var localIP string

func LocalIP() (string, error) {
	if localIP == "" {
		conn, err := net.DialTCP("tcp4", nil, &net.TCPAddr{IP: []byte{8, 8, 8, 8}, Port: 53})
		if err != nil {
			return "", err
		}
		defer conn.Close()
		localIP = strings.Split(conn.LocalAddr().String(), ":")[0]
	}
	return localIP, nil
}

func GetPeerID() string {
	localIP, err := LocalIP()
	if err != nil {
		fmt.Println(err)
		localIP = "DISElevatorCONNECTED"
	}
	return fmt.Sprintf("peer-%s", localIP)
}