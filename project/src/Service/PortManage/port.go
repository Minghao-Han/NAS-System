package PortManage

import (
	"context"
	"fmt"
	"math/rand"
	"nas/project/src/Utils"
	"net"
	"sync"
	"time"
)

type Port struct {
	csPort              int
	dsPort              int
	activeConnections   []Connection
	totalConnection     int
	totalConnectionLock chan bool
}

func NewPort(csPort int, dsPort int) *Port {
	port := &Port{
		csPort:              csPort,
		dsPort:              dsPort,
		activeConnections:   make([]Connection, connPerPort),
		totalConnection:     0,
		totalConnectionLock: make(chan bool, connPerPort),
	}
	for index, _ := range port.activeConnections {
		port.activeConnections[index].Initialize()
	}
	port.totalConnectionLock <- true
	return port
}

func (port *Port) GetCsPort() int {
	return port.csPort
}

func (port *Port) GetDsPort() int {
	return port.dsPort
}
func (port *Port) FindConnection(sourceIP net.IP) (*Connection, bool) {
	for index, conn := range port.activeConnections { //有预留
		if conn.GetSourceIP().Equal(sourceIP) {
			return &(port.activeConnections[index]), true
		}
	}
	//无预留或预留已过期
	return nil, false
}

func (port *Port) PrepareNewConnection(portIndex int, sourceIP net.IP, ctx *context.Context, wg *sync.WaitGroup, found chan *Port, connIndexChan chan int) {
	//found: 若本port可以接收这个连接则found<-本端口
	defer func() {
		if err := recover(); err != nil {
			//fmt.Println(err) //处理写入found,connIndexChan可能导致的写入关闭channel的异常
			wg.Done()
		} else {
			wg.Done()
		}
	}()
	<-port.totalConnectionLock
	if port.totalConnection >= connPerPort {
		fmt.Println(port.totalConnection)
		port.totalConnectionLock <- true
		return
	}
	port.totalConnection += 1        //先占用
	port.totalConnectionLock <- true //解锁
	for _, conn := range port.activeConnections {
		select {
		case <-(*ctx).Done(): //监听其他协程是否结束
			return
		default:
			if conn.GetSourceIP().Equal(sourceIP) { //该端口有相同ip的连接，则不能连接该端口
				port.totalConnection -= 1 //解除占用
				return
			}
		}
	}
	randomSource := rand.NewSource(time.Now().UnixNano())           // 创建随机数源
	randomGenerator := rand.New(randomSource)                       // 创建随机数生成器
	randomSleep := randomGenerator.Intn(1000) * min(portIndex, 400) // 生成 1 到 100 范围内的随机整数
	time.Sleep(time.Duration(randomSleep) * time.Nanosecond)
	select {
	case <-(*ctx).Done():
		port.totalConnection -= 1
		return
	default:
		break
	}
	//本port可以接收这个连接，预留connection
	connIndex, _ := port.reserveConnection(sourceIP)
	connIndexChan <- connIndex
	close(connIndexChan)
	found <- port
	go func() { //处理超时无连接
		timeoutChan := make(chan bool, 1)              //超时通道
		Utils.MakeTimeout(timeoutChan, 5, time.Minute) //
		select {
		case <-timeoutChan:
			//超时未连接就将预留取消
			port.DisConnectByIndex(connIndex, sourceIP)
		case <-port.activeConnections[connIndex].ds2cs:
			return
		}
	}()
	return
}

func (port *Port) reserveConnection(sourceIP net.IP) (int, bool) { //reserveConnection 只能由PrepareNewConnection调用，所以设置为private
	ctx, cancel := context.WithCancel(context.Background())
	set := make(chan int, 1)
	var wg = sync.WaitGroup{}
	wg.Add(len(port.activeConnections))
	for index, _ := range port.activeConnections {
		index := index
		conn := &(port.activeConnections[index])
		go func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err) //处理写入set可能导致的写入关闭channel的异常
					/**
					发生死锁
					*/
					conn.SetSourceIP(nil) //如果发生了写入关闭的chan则说明已经有其他协程找到可用的conn，所以把本协程的预留消除
					wg.Done()
				} else {
					wg.Done()
				}
			}()
			select {
			case <-ctx.Done(): //监听其他协程是否结束
				return
			default:
				conn.LockIP() //获取写锁后也要判断其他协程是否已找到nilConnection，因为获取写锁可能需要一段时间(被其他读写锁阻塞)
				select {
				case <-ctx.Done(): //监听其他协程是否结束
					return
				default:
					//如果是nil
					if conn.sourceIp.Equal(nil) { //这里用conn.sourceIp而不是get方法是因为刚用lock，用get会申请rlock，导致死锁
						conn.sourceIp = sourceIP
						conn.UnlockIP()
						set <- index
					} else {
						conn.UnlockIP() //解锁不能到defer再执行，因为第一个select没有获取写锁
					}
				}
			}
			return
		}()
	}
	//allFinished接收全部完成消息
	allFinished := make(chan bool, 1)
	go func() {
		wg.Wait()
		allFinished <- true
	}()

	select {
	case connIndex := <-set: //找到合适的端口
		close(set)
		cancel()
		wg.Wait()
		return connIndex, true
	case <-allFinished: //全部端口都查过但是没找到合适的
		cancel()
		close(set)
		break
	}
	close(allFinished)
	return -1, false
}

func (port *Port) DisConnectByIP(sourceIP net.IP) bool {
	for index, _ := range port.activeConnections {
		conn := port.activeConnections[index]
		if conn.GetSourceIP().Equal(sourceIP) {
			port.activeConnections[index].Reset()
		}
	}
	return false
}

func (port *Port) DisConnectByIndex(connIndex int, sourceIP net.IP) bool {
	if ok := port.activeConnections[connIndex].SetToIfEqual(sourceIP, nil); ok {
		port.totalConnection -= 1
		fmt.Printf("connection %d on port %d/%d is expired\n", connIndex, port.csPort, port.dsPort)
		return true
	}
	return false
}
