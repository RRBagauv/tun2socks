package mirror

import (
	"github.com/RRBagauv/tun2socks/core/adapter"
	"github.com/RRBagauv/tun2socks/tunnel"
)

var _ adapter.TransportHandler = (*Tunnel)(nil)

type Tunnel struct{}

func (*Tunnel) HandleTCP(conn adapter.TCPConn) {
	tunnel.TCPIn() <- conn
}

func (*Tunnel) HandleUDP(conn adapter.UDPConn) {
	tunnel.UDPIn() <- conn
}
