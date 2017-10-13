package tofino

import (
	"encoding/binary"
	"github.com/bocon13/sdn-go/p4runtime"
	"github.com/bocon13/sdn-go/proto/p4/tmp"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"io/ioutil"
)

const uint32Bytes = 4

type tofinoConfig struct {
	Name            string
	TofinoFilename  string
	ContextFilename string
}

func New(name, tofinoFilename, cxtFilename string) p4runtime.PipelineConfig {
	return &tofinoConfig{
		name, tofinoFilename, cxtFilename,
	}
}

func (c tofinoConfig) Get() (p4runtime.P4DeviceConfig, error) {
	tofinoBin, err := ioutil.ReadFile(c.TofinoFilename)
	if err != nil {
		return nil, errors.Wrapf(err, "Error reading tofino.bin: %s", c.TofinoFilename)
	}
	cxtJson, err := ioutil.ReadFile(c.ContextFilename)
	if err != nil {
		return nil, errors.Wrapf(err, "Error reading context.json: %s", c.ContextFilename)
	}
	// Build the device config
	data := make([]byte, len(c.Name)+len(tofinoBin)+len(cxtJson)+3*uint32Bytes)
	i := 0
	binary.LittleEndian.PutUint32(data[i:i+uint32Bytes], uint32(len(c.Name)))
	i += uint32Bytes
	i += copy(data[i:], []byte(c.Name))
	binary.LittleEndian.PutUint32(data[i:i+uint32Bytes], uint32(len(tofinoBin)))
	i += uint32Bytes
	i += copy(data[i:], tofinoBin)
	binary.LittleEndian.PutUint32(data[i:i+uint32Bytes], uint32(len(cxtJson)))
	i += uint32Bytes
	i += copy(data[i:], cxtJson)
	// quick assertion to make sure the data array matches our expectations
	if i != len(c.Name)+len(tofinoBin)+len(cxtJson)+3*uint32Bytes {
		return nil, errors.Errorf("unable to build config; wrong size")
	}
	deviceConfig := &p4_tmp.P4DeviceConfig{
		DeviceData: data,
	}
	data, err = proto.Marshal(deviceConfig)
	if err != nil {
		return nil, errors.Wrap(err, "error marshalling device config")
	}
	return data, nil
}
