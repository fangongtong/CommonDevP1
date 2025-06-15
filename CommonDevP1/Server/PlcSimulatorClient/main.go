package main

import (
	"CommonDevP1/Server/PlcSimulatorClient/bs"
)

func main() {

	//bs.CommuHost.StartCommu()
	go bs.CommuHost.StartCommu("0.0.0.0:8099", "127.0.0.1:8098")

	bs.RunState()
	//host := Udp.NewUdpClient("127.0.0.1:8098", "0.0.0.0:8099")
}
