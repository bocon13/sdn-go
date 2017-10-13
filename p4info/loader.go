package p4info

import (
	"github.com/bocon13/sdn-go/proto/p4/config"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"io/ioutil"
)

type P4Info struct {
	//TODO try to collapse this
	Info *p4_config.P4Info
}

func LoadInfo(filename string) (*P4Info, error) {
	i := P4Info{
		Info: &p4_config.P4Info{},
	}
	if d, err := ioutil.ReadFile(filename); err != nil {
		return nil, errors.Wrapf(err, "error reading p4info file: %s", filename)
	} else if err := proto.UnmarshalText(string(d), i.Info); err != nil {
		return nil, errors.Wrapf(err, "error parsing p4info file: %s", filename)
	}
	//return &P4Info{i}, nil
	return &i, nil
}
