package portManage

import "net"

type Connection struct {
	sourceIp     net.IP
	cs_on        chan bool //cs stand for control stream
	ds_on        chan bool //ds stand for data stream
	sourceIPLock chan bool
}

func NewConnection() *Connection {
	return &Connection{
		sourceIp:     net.ParseIP("0.0.0.0"),
		cs_on:        make(chan bool, 1),
		ds_on:        make(chan bool, 1),
		sourceIPLock: make(chan bool, 1),
	}
}

func (conn *Connection) SetCSOn() {
	conn.cs_on <- true
}

func (conn *Connection) setDSOn() {
	conn.ds_on <- true
}

func (conn *Connection) SetCSOff() {
	<-conn.cs_on
}

func (conn *Connection) setDSOff() {
	<-conn.ds_on
}

func (conn *Connection) setSourceIP(srcIP net.IP) {
	conn.sourceIp = srcIP
}

func (conn Connection) getSourceIP() net.IP {
	return conn.sourceIp
}
