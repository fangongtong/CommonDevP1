package Trans

import (
	"net"
)

type ITrans interface {
	Send(data []byte) error
	SendTo(data []byte, rUdp *net.UDPAddr) error
	Recv(dt []byte) (int, error)
}
