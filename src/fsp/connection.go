package fsp

import "net"

type Connection struct {
	sourceIp net.IPAddr
	cs_on    bool //cs stand for control stream
	ds_on    bool //ds stand for data stream
}

func (conn *Connection) SetCSOn() {
	conn.cs_on = true
}

func (conn *Connection) setDSOn() {
	conn.ds_on = true
}

func (conn *Connection) getCSOn() bool {
	return conn.cs_on
}

func (conn *Connection) getDSOn() bool {
	return conn.ds_on
}

func (conn *Connection) SetCSOff() {
	conn.cs_on = false
}

func (conn *Connection) setDSOff() {
	conn.ds_on = false
}

func (conn *Connection) setSourceIP(srcIP net.IPAddr) {
	conn.sourceIp = srcIP
}
