package main

import (
	"fmt"
	"nas/project/src/PortManage"
	"net"
	"time"
)

func main() {
	portManage := PortManage.NewPortsManager()
	for i := 0; i < 2; i++ {
		go func() {
			csPort, dsPort, connIndex, ok := portManage.PrepareConnection(net.ParseIP("192.168.1.1"))
			fmt.Printf("csPort is %d,dsPort is %d,connIndex is %d,ok is %t \n", csPort, dsPort, connIndex, ok)
		}()
	}
	time.Sleep(5 * time.Second)
	return
}
