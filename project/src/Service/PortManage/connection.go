package PortManage

import (
	"net"
	"sync"
)

type Connection struct {
	sourceIp net.IP
	cs2ds    chan bool //cs stands for control stream
	ds2cs    chan bool //ds stands for data stream
	ipRWLock *sync.RWMutex
}

func (conn *Connection) Initialize() {
	conn.sourceIp = nil
	conn.cs2ds = make(chan bool, 1)
	conn.ds2cs = make(chan bool, 1)
	conn.ipRWLock = &sync.RWMutex{}
	return
}

func (conn *Connection) SetCSOn() {
	conn.cs2ds <- true
}

func (conn *Connection) SetDSOn() {
	conn.ds2cs <- true
}

func (conn *Connection) SetCSOff() {
	<-conn.cs2ds
}

func (conn *Connection) SetDSOff() {
	<-conn.ds2cs
}

func (conn *Connection) SetSourceIP(srcIP net.IP) {
	conn.ipRWLock.Lock()
	conn.sourceIp = srcIP
	conn.ipRWLock.Unlock()
}

func (conn Connection) GetSourceIP() net.IP {
	defer func() { conn.ipRWLock.RUnlock() }()
	conn.ipRWLock.RLock()
	return conn.sourceIp
}

func (conn *Connection) SetToIfEqual(compare net.IP, to net.IP) bool { //if equal to "compare" than change to "to"
	defer func() { conn.ipRWLock.Unlock() }()
	conn.ipRWLock.Lock()
	if conn.sourceIp.Equal(compare) {
		conn.sourceIp = to
		return true
	}
	return false
}

func (conn *Connection) Reset() {
	conn.ipRWLock.Lock()
	conn.sourceIp = nil
	conn.ipRWLock.Unlock()
}

func (conn *Connection) LockIP() {
	conn.ipRWLock.Lock()
}

func (conn *Connection) UnlockIP() {
	conn.ipRWLock.Unlock()
}
