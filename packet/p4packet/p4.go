package p4packet

import (
	"fmt"
	"github.com/bocon13/sdn-go/proto/p4"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/pkg/errors"
	"io"
)

var serializeOptions gopacket.SerializeOptions = gopacket.SerializeOptions{
	FixLengths:       true,
	ComputeChecksums: true,
}

type Listener func(pkt gopacket.Packet, metadata []*p4.PacketMetadata)

func ListenViaP4(stream p4.P4Runtime_StreamChannelClient, listener Listener) chan interface{} {
	stop := make(chan interface{})
	go func() {
		fmt.Println("Listening for pkt ins...")
		for {
			//FIXME the stop channel currently does nothing
			//_, ok := <-stop
			//if !ok {
			//	return
			//}
			resp, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("eof")
				return
			}
			if err != nil {
				fmt.Println("Error receiving packet on stream", err)
				continue
			}
			switch t := resp.Update.(type) {
			case *p4.StreamMessageResponse_Arbitration:
				fmt.Println("arb:", t)
			case *p4.StreamMessageResponse_Packet:
				if t.Packet != nil {
					//fmt.Println(hex.Dump(t.Packet.Payload))
					pkt := gopacket.NewPacket(t.Packet.Payload, layers.LayerTypeEthernet, gopacket.Default)
					fmt.Printf("Packet received via P4Runtime:\n%s\n", pkt)
					if listener != nil {
						listener(pkt, t.Packet.Metadata)
					}
				}
			default:
				fmt.Println("unknown", resp)
			}
		}
		fmt.Println("end")
	}()
	return stop
}

func SendViaP4(stream p4.P4Runtime_StreamChannelClient, pkt gopacket.Packet, metadata []*p4.PacketMetadata) error {
	buf := gopacket.NewSerializeBuffer()
	gopacket.SerializePacket(buf, serializeOptions, pkt)
	err := SendViaP4Raw(stream, buf.Bytes(), metadata)
	//if err == nil {
	//	fmt.Printf("Packet sent via P4Runtime:\n%s\n", pkt)
	//}
	return err
}

func SendViaP4Raw(stream p4.P4Runtime_StreamChannelClient, data []byte, metadata []*p4.PacketMetadata) error {
	err := stream.Send(&p4.StreamMessageRequest{
		Update: &p4.StreamMessageRequest_Packet{
			Packet: &p4.PacketOut{
				Payload:  data,
				Metadata: metadata,
			},
		},
	})
	if err != nil {
		return errors.Wrap(err, "error sending packet on stream")
	}
	return nil
}
