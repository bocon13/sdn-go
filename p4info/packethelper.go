package p4info

import (
	"fmt"
	"github.com/bocon13/sdn-go/proto/p4/config"
)

const PACKET_IN_ANNOTATIONS = "packet_in"
const PACKET_OUT_ANNOTATIONS = "packet_out"

type ControllerPacketMetadata struct {
	*p4_config.ControllerPacketMetadata
}

func (info *P4Info) GetPacketAnnotations(name string) (ControllerPacketMetadata, error) {
	for i := range info.Info.ControllerPacketMetadata {
		t := info.Info.ControllerPacketMetadata[i]
		if _, ok := matches(t, name); ok {
			return ControllerPacketMetadata{t}, nil
		}
	}
	return ControllerPacketMetadata{}, fmt.Errorf("controller packet metadata %s not found", name)
}

func (info *P4Info) GetPacketInAnnotations() (ControllerPacketMetadata, error) {
	return info.GetPacketAnnotations(PACKET_IN_ANNOTATIONS)
}

func (info *P4Info) GetPacketOutAnnotations() (ControllerPacketMetadata, error) {
	return info.GetPacketAnnotations(PACKET_OUT_ANNOTATIONS)
}

func (pm *ControllerPacketMetadata) GetMetadata(name string) *p4_config.ControllerPacketMetadata_Metadata {
	for i := range pm.Metadata {
		m := pm.Metadata[i]
		if m.Name == name {
			return m
		}
	}
	return nil //FIXME
}
