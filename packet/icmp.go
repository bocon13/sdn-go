package packet

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"golang.org/x/net/ipv4"
	"math/rand"
	"net"
)

const DEFAULT_TTL uint8 = 64
const DEFAULT_PING_PAYLOAD string = "abcdefghijklmnopqrstuvwxyz"

var serializeOptions gopacket.SerializeOptions = gopacket.SerializeOptions{
	FixLengths:       true,
	ComputeChecksums: true,
}

func GenICMPEchoReq(srcMac, dstMac net.HardwareAddr, srcIp, dstIp net.IP, id, seqNo uint16) []byte {
	buf := gopacket.NewSerializeBuffer()
	gopacket.SerializeLayers(buf, serializeOptions,
		&layers.Ethernet{
			SrcMAC:       srcMac,
			DstMAC:       dstMac,
			EthernetType: layers.EthernetTypeIPv4,
		},
		&layers.IPv4{
			Version:  ipv4.Version,
			Protocol: layers.IPProtocolICMPv4,
			TTL:      DEFAULT_TTL,
			Id:       uint16(rand.Int()), //TODO should probably be smarter about this
			SrcIP:    srcIp,
			DstIP:    dstIp,
		},
		&layers.ICMPv4{
			TypeCode: layers.ICMPv4TypeEchoRequest << 8,
			Id:       id,
			Seq:      seqNo,
		},
		gopacket.Payload([]byte(DEFAULT_PING_PAYLOAD)),
	)
	return buf.Bytes()
}
