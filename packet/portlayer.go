package packet

import (
	"encoding/binary"
	"errors"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

/*

	This is just an example of how to create new gopacket layers.

	This layer consists of a 2 byte "port" field and encapsulates Ethernet frames.

*/

// Register the port layer, giving it name "Port" and its decoder
// Note: num must be unique across all decoders
var LayerTypePort = gopacket.RegisterLayerType(1001,
	gopacket.LayerTypeMetadata{Name: "Port", Decoder: gopacket.DecodeFunc(decodePort)})

// Port is the custom P4 frame with ingress/egress port number
type Port struct {
	layers.BaseLayer
	Port            uint16
	SubmitToIngress bool
}

func (p *Port) LayerType() gopacket.LayerType { return LayerTypePort }

func decodePort(data []byte, p gopacket.PacketBuilder) error {
	port := &Port{}
	err := port.DecodeFromBytes(data, p)
	if err != nil {
		return err
	}
	p.AddLayer(port)
	// We assume that the payload is an Ethernet packet
	return p.NextDecoder(layers.LayerTypeEthernet)
}

func (p *Port) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	if len(data) < 2 {
		return errors.New("Packet too small")
	}
	num := binary.BigEndian.Uint16(data[0:2])
	p.Port = num >> 7
	p.SubmitToIngress = (num>>6)&1 == 1
	p.BaseLayer = layers.BaseLayer{data[:2], data[2:]}
	return nil
}

func (p *Port) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	bytes, err := b.PrependBytes(2)
	if err != nil {
		return err
	}
	num := p.Port << 7
	if p.SubmitToIngress {
		num |= 1 << 6
	}
	binary.BigEndian.PutUint16(bytes, num)
	return nil
}
