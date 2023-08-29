package portManage

import (
	"context"
	"fmt"
	"nas/src/myTime"
	"net"
	"sync"
	"time"
)

var (
	portNum                  int
	connPerPort              int
	prepareConnectionTimeout int //use ms as unit
)

type PortsManager struct {
	availableConnections chan bool
	ports                []Port
}

func NewPortsManager() *PortsManager {
	var portsManager = PortsManager{}
	/**
	initial operation, read from configuration file
	*/
	return &portsManager
}

func (pm *PortsManager) FindPort(csPort int, dsPort int) (*Port, bool) { //用csPort查时dsPort设为0
	if dsPort == 0 {
		for _, port := range pm.ports {
			if port.GetCsPort() == csPort {
				return &port, true
			}
		}
	} else if csPort == 0 {
		for _, port := range pm.ports {
			if port.GetDsPort() == dsPort {
				return &port, true
			}
		}
	}
	return nil, false
}

func (pm *PortsManager) PrepareConnection(sourceIP net.IP) (int, int, bool) {
	timeoutChan := make(chan bool, 1)
	myTime.MakeTimeout(timeoutChan, prepareConnectionTimeout, time.Millisecond) //设定超时，由配置文件决定
	select {
	case <-pm.availableConnections: //通过availableConnections实现pv，当总连接数不够时阻塞
		//能连接
		{
			ctx, cancel := context.WithCancel(context.Background())
			found := make(chan *Port, 1)
			var wg = sync.WaitGroup{}
			wg.Add(portNum)
			for _, port := range pm.ports {
				go port.PrepareNewConnection(sourceIP, &ctx, &wg, found)
			}
			allFinished := make(chan bool, 1)
			go func() {
				wg.Wait()
				allFinished <- true
			}()
			select {
			case availablePort := <-found: //找到合适的端口
				cancel()
				return availablePort.GetCsPort(), availablePort.GetDsPort(), true
			case <-allFinished: //全部端口都查过但是没找到合适的
				cancel()
				break
			}
			wg.Wait()
			close(found)
			close(allFinished)
		}
	case <-timeoutChan: //超时退出
		break
	default:
		fmt.Println("prepareConnection wrong")
	}
	return 0, 0, false
}
