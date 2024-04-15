package main

import (
	"fmt"
	"net"
	"net/netip"
	"os"
	"strings"
	"time"

	"github.com/mdlayher/ndp"
)

func main() {
	ifname := os.Args[1]
	//fmt.Printf("interface: %s\n", ifname)

	addr, err := sendRS(ifname)
	if err != nil {
		panic(err)
	}

	a := strings.Split(addr, "%")[0]
	fmt.Printf("define %s_neighbor = %s;\n", ifname, a)
}

func sendRS(ifname string) (string, error) {
	iface, err := net.InterfaceByName(ifname)
	if err != nil {
		return "", err
	}

	c, _, err := ndp.Listen(iface, ndp.Addr("linklocal"))
	//c, ip, err := ndp.Listen(iface, ndp.Addr("linklocal"))
	if err != nil {
		return "", err
	}

	hwaddr := iface.HardwareAddr
	//fmt.Printf("ip: %s, source link-layer address: %s", ip, hwaddr)

	m := &ndp.RouterSolicitation{}
	m.Options = append(m.Options, &ndp.LinkLayerAddress{
		Direction: ndp.Source,
		Addr:      hwaddr,
	})

	for {
		dst := netip.MustParseAddr("ff02::2")
		err = c.WriteTo(m, nil, dst)
		if err != nil {
			return "", err
		}
		err = c.SetReadDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			return "", err
		}

		msg, _, from, err := c.ReadFrom()
		if err == nil {
			_, ok := msg.(*ndp.RouterAdvertisement)
			if ok {
				return from.String(), nil
			}
		}

		fmt.Printf(".")
		time.Sleep(1 * time.Second)
	}
}
