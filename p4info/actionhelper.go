package p4info

import (
	"fmt"
	"github.com/bocon13/sdn-go/proto/p4"
	"github.com/bocon13/sdn-go/proto/p4/config"
)

type Action struct {
	info *P4Info
	Id   uint32
	Name string
	*p4_config.Action
}

func (info *P4Info) GetAction(name string) (Action, error) {
	for i := range info.Info.Actions {
		a := info.Info.Actions[i]
		if p, ok := matches(a, name); ok {
			return Action{info, p.Id, name, a}, nil
		}
	}
	return Action{}, fmt.Errorf("Action %s not found", name)
}

var emptyMap map[string][]byte = make(map[string][]byte)

func (a *Action) GetTableAction(params map[string][]byte) (*p4.TableAction, error) {
	if params != nil {
		//TODO this is kind of a hack, so that we don't have to do nil checking everywhere
		params = emptyMap
	}

	usedParams := make(map[string]bool)
	for k := range params {
		usedParams[k] = false
	}

	p4Params := make([]*p4.Action_Param, len(a.GetParams()))
	for i := range a.GetParams() {
		p := a.GetParams()[i]
		ap := &p4.Action_Param{
			ParamId: p.Id,
		}
		bytes := RoundedByte(p.Bitwidth)
		if v, ok := params[p.Name]; !ok {
			// param not set, use default value
			ap.Value = make([]byte, bytes)
		} else if len(v) != int(bytes) {
			// param value differs from required value
			return nil, fmt.Errorf("value for param %s is %d bytes (required %d bytes)", p.Name, len(v), bytes)
		} else {
			// using supplied value
			ap.Value = v
			usedParams[p.Name] = true
		}
		p4Params[i] = ap
	}

	errStr := ""
	for k, v := range usedParams {
		if !v {
			errStr = k + ", "
		}
	}
	if errStr != "" {
		return nil, fmt.Errorf("action %s has extra parameters: %s", a.Name, errStr)
	}

	return &p4.TableAction{
		Type: &p4.TableAction_Action{
			Action: &p4.Action{
				ActionId: a.Id,
				Params:   p4Params,
			},
		},
	}, nil
}

func (a *Action) GetMemberAction(id uint32) *p4.TableAction {
	return &p4.TableAction{
		Type: &p4.TableAction_ActionProfileMemberId{
			ActionProfileMemberId: id,
		},
	}
}

func (a *Action) GetGroupAction(id uint32) *p4.TableAction {
	return &p4.TableAction{
		Type: &p4.TableAction_ActionProfileGroupId{
			ActionProfileGroupId: id,
		},
	}
}
