package p4runtime

import (
	"context"
	"fmt"
	"github.com/bocon13/sdn-go/proto/p4"
	"github.com/bocon13/sdn-go/proto/p4/config"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type P4DeviceConfig []byte

type PipelineConfig interface {
	Get() (P4DeviceConfig, error)
}

func GetPipelineConfigs(client p4.P4RuntimeClient, deviceId uint64) (*p4.ForwardingPipelineConfig, error) {
	getReq := &p4.GetForwardingPipelineConfigRequest{
		DeviceIds: []uint64{deviceId},
	}
	replies, err := client.GetForwardingPipelineConfig(context.Background(), getReq)
	//TODO update ErrorDesc to use non-deprecated method
	if grpc.ErrorDesc(err) == "No forwarding pipeline config set for this device" {
		fmt.Println("no forwarding pipeline; need to push one")
		return &p4.ForwardingPipelineConfig{}, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "error getting pipeline config")
	}
	for i := range replies.Configs {
		c := replies.Configs[i]
		if c.DeviceId == deviceId {
			return c, nil
		}
	}
	return &p4.ForwardingPipelineConfig{}, nil
}

func SetPipelineConfig(client p4.P4RuntimeClient, config *p4.ForwardingPipelineConfig) error {
	// Need to push the config
	setReq := p4.SetForwardingPipelineConfigRequest{
		Action:  p4.SetForwardingPipelineConfigRequest_VERIFY_AND_COMMIT,
		Configs: []*p4.ForwardingPipelineConfig{config},
	}
	_, err := client.SetForwardingPipelineConfig(context.Background(), &setReq)
	return errors.Wrap(err, "error setting pipeline config")
}

func matches(target, actual *p4.ForwardingPipelineConfig) bool {
	// TODO Tofino doesn't appear to fill in the device config on Get, so assume it matches
	// When it does, we can replace this with proto compare: proto.Equal(target, actual)
	return target.DeviceId == actual.DeviceId && proto.Equal(target.P4Info, actual.P4Info)
}

func UpdatePipelineConfig(client p4.P4RuntimeClient, p4Info *p4_config.P4Info,
	config PipelineConfig, deviceId uint64, forcePush bool) (bool, error) {
	configData, err := config.Get()
	if err != nil {
		return false, errors.Wrap(err, "error building target config")
	}
	targetConfig := &p4.ForwardingPipelineConfig{
		DeviceId:       deviceId,
		P4Info:         p4Info,
		P4DeviceConfig: configData,
	}

	deviceConfig, err := GetPipelineConfigs(client, deviceId)
	if err != nil {
		return false, errors.Wrap(err, "error getting device config")
	}

	if forcePush || !matches(targetConfig, deviceConfig) {
		// Config doesn't match or updated is forced, so re-push...
		err = SetPipelineConfig(client, targetConfig)
		if err != nil {
			return true, errors.Wrap(err, "error setting config")
		}
		return true, nil
	}
	return false, nil
}
