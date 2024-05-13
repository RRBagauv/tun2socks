package tunnel

import (
	"github.com/RRBagauv/tun2socks/core/adapter"
	"gvisor.dev/gvisor/pkg/log"
)

// Unbuffered TCP/UDP queues.
var (
	_tcpQueue = make(chan adapter.TCPConn)
	_udpQueue = make(chan adapter.UDPConn)
)

func init() {
	log.Infof(
		"Start process",
	)
	go process()
}

// TCPIn return fan-in TCP queue.
func TCPIn() chan<- adapter.TCPConn {
	return _tcpQueue
}

// UDPIn return fan-in UDP queue.
func UDPIn() chan<- adapter.UDPConn {
	return _udpQueue
}

func process() {
	for {
		select {
		case conn := <-_tcpQueue:
			go handleTCPConn(conn)
		case conn := <-_udpQueue:
			go handleUDPConn(conn)
		}
	}
}
