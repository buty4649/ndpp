/*
Copyright © 2024 buty4649

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package ndp

import (
	"context"
	"net"
	"net/netip"
	"time"

	"github.com/mdlayher/ndp"
)

type Result struct {
	Router NetAddr
	Local  NetAddr
	IfName string
}

type NetAddr struct {
	Addr   netip.Addr
	LLAddr string
}

func SendRS(ctx context.Context, ifname string) (*Result, error) {
	iface, err := net.InterfaceByName(ifname)
	if err != nil {
		return nil, err
	}

	c, localAddr, err := ndp.Listen(iface, ndp.Addr("linklocal"))
	if err != nil {
		return nil, err
	}

	hwaddr := iface.HardwareAddr

	m := &ndp.RouterSolicitation{}
	m.Options = append(m.Options, &ndp.LinkLayerAddress{
		Direction: ndp.Source,
		Addr:      hwaddr,
	})

	dst := netip.MustParseAddr("ff02::2")
	for {
		err = c.WriteTo(m, nil, dst)
		if err != nil {
			return nil, err
		}
		err = c.SetReadDeadline(time.Now().Add(1 * time.Second))
		if err != nil {
			return nil, err
		}

		msg, _, from, err := c.ReadFrom()
		if err == nil {
			r, ok := msg.(*ndp.RouterAdvertisement)
			if ok {
				router_addr := netip.MustParseAddr(from.String())
				router_lladdr := extractLinkLayerAddress(r)
				lladdr := iface.HardwareAddr.String()
				return &Result{
					Router: NetAddr{
						Addr:   router_addr,
						LLAddr: router_lladdr,
					},
					Local: NetAddr{
						Addr:   localAddr,
						LLAddr: lladdr,
					},
					IfName: ifname,
				}, nil
			}
		}

		if nerr, ok := err.(net.Error); ok && !nerr.Timeout() {
			return nil, err
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
	}
}

func extractLinkLayerAddress(r *ndp.RouterAdvertisement) string {
	o := r.Options

	if len(o) == 0 {
		return ""
	}

	return o[0].(*ndp.LinkLayerAddress).Addr.String()
}
