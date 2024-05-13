package tunnel

import (
	"io"
	"net"
	"sync"
	"time"

	"github.com/RRBagauv/tun2socks/common/pool"
	"github.com/RRBagauv/tun2socks/core/adapter"
	"github.com/RRBagauv/tun2socks/log"
	M "github.com/RRBagauv/tun2socks/metadata"
	"github.com/RRBagauv/tun2socks/proxy"
	"github.com/RRBagauv/tun2socks/tunnel/statistic"
)

const (
	// tcpWaitTimeout implements a TCP half-close timeout.
	tcpWaitTimeout = 60 * time.Second
)

func handleTCPConn(originConn adapter.TCPConn) {
	defer originConn.Close()

	id := originConn.ID()
	metadata := &M.Metadata{
		Network: M.TCP,
		SrcIP:   net.IP(id.RemoteAddress.AsSlice()),
		SrcPort: id.RemotePort,
		DstIP:   net.IP(id.LocalAddress.AsSlice()),
		DstPort: id.LocalPort,
	}

	remoteConn, err := proxy.Dial(metadata)
	if err != nil {
		log.Warnf("[TCP] dial %s: %v", metadata.DestinationAddress(), err)
		return
	}
	metadata.MidIP, metadata.MidPort = parseAddr(remoteConn.LocalAddr())

	remoteConn = statistic.DefaultTCPTracker(remoteConn, metadata)
	defer remoteConn.Close()

	log.Infof("[TCP] %s <-> %s", metadata.SourceAddress(), metadata.DestinationAddress())
	pipe(originConn, remoteConn)
}

// pipe copies copy data to & from provided net.Conn(s) bidirectionally.
func pipe(origin, remote net.Conn) {
	wg := sync.WaitGroup{}
	wg.Add(2)

	go unidirectionalStream(remote, origin, "origin->remote", &wg)
	go unidirectionalStream(origin, remote, "remote->origin", &wg)

	wg.Wait()
}

func unidirectionalStream(dst, src net.Conn, dir string, wg *sync.WaitGroup) {
	defer wg.Done()
	buf := pool.Get(pool.RelayBufferSize)
	if _, err := io.CopyBuffer(dst, src, buf); err != nil {
		log.Debugf("[TCP] copy data for %s: %v", dir, err)
	}
	pool.Put(buf)
	// Do the upload/download side TCP half-close.
	if cr, ok := src.(interface{ CloseRead() error }); ok {
		cr.CloseRead()
	}
	if cw, ok := dst.(interface{ CloseWrite() error }); ok {
		cw.CloseWrite()
	}
	// Set TCP half-close timeout.
	dst.SetReadDeadline(time.Now().Add(tcpWaitTimeout))
}
