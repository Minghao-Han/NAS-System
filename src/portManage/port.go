package portManage

import (
	"context"
	"fmt"
	"nas/src/myTime"
	"net"
	"sync"
	"time"
)

type Port struct {
	csPort              int
	dsPort              int
	connections         []Connection
	totalConnection     int
	totalConnectionLock chan bool
}

func NewPort(csPort int, dsPort int) *Port {
	return &Port{
		csPort:              csPort,
		dsPort:              dsPort,
		connections:         make([]Connection, connPerPort),
		totalConnection:     0,
		totalConnectionLock: make(chan bool, connPerPort)}
}

func (port *Port) GetCsPort() int {
	return port.csPort
}

func (port *Port) GetDsPort() int {
	return port.dsPort
}
func (port *Port) FindConnection(sourceIP net.IP) (*Connection, bool) {
	for _, conn := range port.connections { //有预留
		if conn.sourceIp.Equal(sourceIP) {
			return &conn, true
		}
	}
	//无预留或预留已过期
	return nil, false
}

func (port *Port) PrepareNewConnection(sourceIP net.IP, ctx *context.Context, wg *sync.WaitGroup, found chan *Port) {
	//found: 若本port可以接收这个连接则found<-本端口
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			wg.Done()
		}
		wg.Done()
	}()
	<-port.totalConnectionLock
	if port.totalConnection >= connPerPort {
		port.totalConnectionLock <- true
		return
	}
	port.totalConnection += 1        //先占用
	port.totalConnectionLock <- true //解锁
	duplicateIP := false
	connIndex := 0
	for i, conn := range port.connections { //这里可能需要改
		select {
		case <-(*ctx).Done(): //监听其他协程是否结束
			return
		default:
			if conn.sourceIp.Equal(sourceIP) { //该端口有相同ip的连接
				duplicateIP = true
			} else if conn.sourceIp.Equal(net.ParseIP("0.0.0.0")) { //0.0.0.0的连接位置，即available
				connIndex = i //记录可用conn
			}
		}
	}
	if duplicateIP {
		port.totalConnection -= 1 //解除占用
		return
	} else {
		//本port可以接收这个连接，预留connection
		port.connections[connIndex].sourceIp = sourceIP
		found <- port
		go func() { //处理超时无连接
			timeoutChan := make(chan bool, 1) //超时通道
			myTime.MakeTimeout(timeoutChan, 500, time.Millisecond)
			select {
			case <-timeoutChan:
				{ //超时未连接就将预留取消
					port.totalConnection -= 1 //解除占用
					port.connections[connIndex].sourceIp = net.ParseIP("0.0.0.0")
				}
			case <-port.connections[connIndex].cs_on: //控制流连接成功
				select {
				case <-timeoutChan: //超时
					port.connections[connIndex].sourceIp = net.ParseIP("0.0.0.0")
					port.totalConnection -= 1 //解除占用
				case <-port.connections[connIndex].ds_on:
					return
				}
			}
		}()
		return
	}
}
