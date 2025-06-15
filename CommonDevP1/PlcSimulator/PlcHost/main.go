package main

import (
	dtf "CommonDevP1/PlcSimulator/DataTransformer"
	dm "CommonDevP1/PlcSimulator/DeviceMemory"
	"CommonDevP1/PlcSimulator/PlcHost/bs"
	"CommonDevP1/PlcSimulator/Trans/Udp"
	"fmt"
)

func main() {

	host := Udp.NewUdpHost("0.0.0.0:8098")
	if host == nil {
		return
	}
	plcBs := bs.NewBs()
	var recvBuf [1048]byte
	for {
		n, e := host.Recv(recvBuf[:])
		if e != nil {
			fmt.Printf("Error on recv: %s \r\n", e.Error())
			return
		}

		code, req := dtf.UnPackData1(recvBuf[:n])

		fmt.Printf("recved code is:%d\r\n", code)

		switch code {
		case 0: //  0:real data request
			realDt := plcBs.GetRealData()
			//  2: real data return
			dt := dtf.PackData(2, realDt)
			host.Send(dt)
		case 1: //  1:cmd
			plcBs.SetCmd(req.(*dm.DBA_1))
		}

	}
}
