package p4info

import (
	"github.com/bocon13/sdn-go/proto/p4/config"
)

type hasPreamble interface {
	GetPreamble() *p4_config.Preamble
}

func matches(e hasPreamble, name string) (*p4_config.Preamble, bool) {
	p := e.GetPreamble()
	return p, p != nil && (p.Name == name || p.Alias == name)
}
