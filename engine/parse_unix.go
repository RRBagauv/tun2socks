//go:build unix

package engine

import (
	"net/url"

	"github.com/RRBagauv/tun2socks/core/device"
	"github.com/RRBagauv/tun2socks/core/device/tun"
)

func parseTUN(u *url.URL, mtu uint32) (device.Device, error) {
	return tun.Open(u.Host, mtu)
}
