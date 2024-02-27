package PortManage

import (
	"context"
	"nas/project/src/Utils"
	"net"
	"sync"
	"time"
)

var (
	portNum                  int
	connPerPort              int
	prepareConnectionTimeout int //use s as unit
)

type PortsManager struct {
	availableConnections chan bool
	ports                []Port
}

var defaultPortsManager *PortsManager

func DefaultPortsManager() *PortsManager {
	if defaultPortsManager == nil {
		/**
		initial operation, read from configuration file
		*/
		csPorts := Utils.DefaultConfigReader().Get("FSP:csPorts").([]interface{})
		dsPorts := Utils.DefaultConfigReader().Get("FSP:dsPorts").([]interface{})
		portNum = len(dsPorts)
		connPerPort = Utils.DefaultConfigReader().Get("FSP:connPerPort").(int)
		prepareConnectionTimeout = Utils.DefaultConfigReader().Get("FSP:prepareConnectionTimeout").(int)
		var portsManager = PortsManager{
			availableConnections: make(chan bool, connPerPort*portNum),
			ports:                make([]Port, 0),
		}
		for i := 0; i < connPerPort*portNum; i++ {
			portsManager.availableConnections <- true
		}
		for index, csPort := range csPorts {
			portsManager.ports = append(portsManager.ports, *NewPort(csPort.(int), dsPorts[index].(int)))
		}
		defaultPortsManager = &portsManager
	}
	return defaultPortsManager
}

// FindPort csPort,dsPort
func (pm *PortsManager) FindPort(csPort int, dsPort int) (*Port, bool) { //用csPort查时dsPort设为0
	if dsPort == 0 {
		for index, _ := range pm.ports {
			port := &(pm.ports[index])
			if port.GetCsPort() == csPort {
				return port, true
			}
		}
	} else if csPort == 0 {
		for index, _ := range pm.ports {
			port := &(pm.ports[index])
			if port.GetDsPort() == dsPort {
				return port, true
			}
		}
	}
	return nil, false
}

// PrepareConnection 返回值分别是，csPort,dsPort,connIndex,成功预留
func (pm *PortsManager) PrepareConnection(sourceIP net.IP) (int, int, int, bool) {
	timeoutChan := make(chan bool, 1)
	go Utils.MakeTimeout(timeoutChan, prepareConnectionTimeout, time.Millisecond) //设定超时，由配置文件决定
	select {
	case <-pm.availableConnections: //通过availableConnections实现pv，当总连接数不够时阻塞
		//能连接
		{
			ctx, cancel := context.WithCancel(context.Background())
			found := make(chan *Port, 1)
			connIndexChan := make(chan int, 1)
			var wg = sync.WaitGroup{}
			wg.Add(portNum)
			ctxLock := &sync.Mutex{}
			for index, _ := range pm.ports {
				port := &pm.ports[index]
				go port.PrepareNewConnection(sourceIP, &ctx, ctxLock, &wg, found, connIndexChan, cancel)
			}
			allFinished := make(chan bool, 1)
			//allFinished接收超时消息
			go func() {
				wg.Wait()
				allFinished <- true
			}()
			select {
			case availablePort := <-found: //找到合适的端口
				close(found)
				connIndex := <-connIndexChan
				wg.Wait()
				return availablePort.GetCsPort(), availablePort.GetDsPort(), connIndex, true
			case <-allFinished: //全部端口都查过但是没找到合适的
				cancel()
				close(found)
				break
			}
			close(allFinished)
		}
	case <-timeoutChan: //超时退出
		break
	}
	return 0, 0, -1, false
}
