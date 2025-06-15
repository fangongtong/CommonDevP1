package Udp

import (
	"CommonDevP1/PlcSimulator/Trans"
	"fmt"
	"net"
	"time"
)

type _Udp struct {
	conn      *net.UDPConn
	remoteUdp *net.UDPAddr
}

func (this *_Udp) Send(data []byte) error {
	//  模仿真实延时,单方10ms,来回正好20ms
	time.Sleep(time.Millisecond * 10)

	_, err := this.conn.WriteToUDP(data, this.remoteUdp)
	return err
}
func (this *_Udp) SendTo(data []byte, rUdp *net.UDPAddr) error {
	//  模仿真实延时,单方10ms,来回正好20ms
	time.Sleep(time.Millisecond * 10)

	_, err := this.conn.WriteToUDP(data, rUdp)
	return err
}

func (this *_Udp) Recv(dt []byte) (int, error) {

	n, rAddr, e := this.conn.ReadFromUDP(dt)

	if e != nil {
		fmt.Println("ReadFromUDP failed")
		return 0, e
	} else {
		fmt.Println("ReadFromUDP ok")
	}

	this.remoteUdp = rAddr

	return n, nil
}

func NewUdpHost(ipPort string) Trans.ITrans {

	var conn *net.UDPConn
	if addr, e := net.ResolveUDPAddr("udp", ipPort); e != nil {
		return nil
	} else if conn, e = net.ListenUDP("udp", addr); e != nil {
		return nil
	}

	fmt.Println("Udp Listen ok")
	return &_Udp{conn: conn}
}

func NewUdpClient(remoteIpPort, localIpPort string) Trans.ITrans {

	var conn *net.UDPConn
	var rAddr *net.UDPAddr
	var e error
	// if lAddr, e := net.ResolveUDPAddr("udp", localIpPort); e != nil {
	// 	return nil
	// } else
	if rAddr, e = net.ResolveUDPAddr("udp", remoteIpPort); e != nil {
		return nil
	} else if conn, e = net.DialUDP("udp", nil, rAddr); e != nil {
		return nil
	}
	return &_Udp{conn: conn, remoteUdp: rAddr}
}
