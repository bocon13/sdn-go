package pcap

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/pkg/errors"
	"net"
	"time"
)

const MAX_PACKET_SIZE int32 = 65536
const DEFAULT_TIMEOUT time.Duration = 0 //5 * time.Second

var serializeOptions gopacket.SerializeOptions = gopacket.SerializeOptions{
	FixLengths:       true,
	ComputeChecksums: true,
}

type Listener func(pkt gopacket.Packet)

// Cache of interface name to pcap handle
var handles map[string]*pcap.Handle = make(map[string]*pcap.Handle)

type IntfInfo struct {
	Name         string
	HardwareAddr net.HardwareAddr
	IP           net.IP
}

func GetInfo(name string) IntfInfo {
	intf, err := net.InterfaceByName(name)
	if err != nil {
		panic(err)
	}
	ip := getIp(intf)
	return IntfInfo{
		Name:         name,
		HardwareAddr: intf.HardwareAddr,
		IP:           ip,
	}
}

func GetHandle(intf string) *pcap.Handle {
	h, ok := handles[intf]
	if ok {
		return h
	}
	h, err := pcap.OpenLive(intf, MAX_PACKET_SIZE, true, DEFAULT_TIMEOUT)
	if err != nil {
		panic(err)
	}
	handles[intf] = h
	return h
}

func ListenViaPcap(intf string, listener Listener) chan interface{} {
	handle := GetHandle(intf)
	stop := make(chan interface{})
	go func() {
		packetSource := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)
		for {
			select {
			case pkt := <-packetSource.Packets():
				if pkt != nil {
					fmt.Printf("Packet received on %s:\n%s\n", intf, pkt)
					if listener != nil {
						listener(pkt)
					}
				}
			case <-stop:
				fmt.Printf("Stopped listening on %s\n", intf)
				return
			}
		}
	}()
	return stop
}

func SendViaPcap(intf string, pkt gopacket.Packet) error {
	buf := gopacket.NewSerializeBuffer()
	gopacket.SerializePacket(buf, serializeOptions, pkt)
	return SendViaPcapRaw(intf, buf.Bytes())
}

func SendViaPcapRaw(intf string, data []byte) error {
	handle := GetHandle(intf)
	err := handle.WritePacketData(data)
	if err != nil {
		return errors.Wrap(err, "error sending packet on stream")
	}
	return nil
}

func CloseHandles() {
	for k, v := range handles {
		v.Close()
		delete(handles, k)
	}
}

func getIp(intf *net.Interface) net.IP {
	addrs, err := intf.Addrs()
	if err != nil {
		panic(err)
	}
	for _, a := range addrs {
		switch t := a.(type) {
		case *net.IPNet:
			return t.IP
		}
	}
	panic("no IP found")
}
