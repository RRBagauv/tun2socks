package fdbased

import (
	"fmt"
	"github.com/RRBagauv/tun2socks/log"
	"os"

	"github.com/RRBagauv/tun2socks/core/device"
	"github.com/RRBagauv/tun2socks/core/device/iobased"
)

func open(fd int, mtu uint32, offset int) (device.Device, error) {
	f := &FD{fd: fd, mtu: mtu}

	ep, err := iobased.New(os.NewFile(uintptr(fd), f.Name()), mtu, 500)
	if err != nil {
		log.Debugf("create endpoint: %w", err)
		return nil, fmt.Errorf("create endpoint: %w", err)
	}
	f.LinkEndpoint = ep
	log.Debugf("Success start")

	return f, nil
}
