package p4info

import (
	"fmt"
	"github.com/bocon13/sdn-go/proto/p4"
	"github.com/bocon13/sdn-go/proto/p4/config"
)

type Table struct {
	info *P4Info
	Id   uint32
	Name string
	*p4_config.Table
}

func (info *P4Info) GetTable(name string) (Table, error) {
	for i := range info.Info.Tables {
		t := info.Info.Tables[i]
		if p, ok := matches(t, name); ok {
			return Table{info, p.Id, name, t}, nil
		}
	}
	return Table{}, fmt.Errorf("table %s not found", name)
}

func (t Table) GetTableEntry(matches map[string]interface{},
	action *p4.TableAction,
	actionParams map[string]*p4.Action_Param,
	priority int32,
	metadata uint64) (*p4.Entity, error) {

	// Build the match fields
	matchFields := make([]*p4.FieldMatch, len(t.MatchFields))
	for i := range t.GetMatchFields() {
		m := t.GetMatchFields()[i]
		customMatch, ok := matches[m.Name]
		if ok {
			fmt.Println(customMatch) //FIXME
		}

		// Build default match
		valueSize := m.Bitwidth / 8
		if m.Bitwidth%8 > 0 {
			valueSize++
		}
		value := make([]byte, valueSize)
		f := &p4.FieldMatch{
			FieldId: m.Id,
		}
		switch m.MatchType {
		case p4_config.MatchField_VALID:
			f.FieldMatchType = &p4.FieldMatch_Valid_{
				Valid: &p4.FieldMatch_Valid{},
			}
		case p4_config.MatchField_EXACT:
			f.FieldMatchType = &p4.FieldMatch_Exact_{
				Exact: &p4.FieldMatch_Exact{
					Value: value,
				},
			}
		case p4_config.MatchField_LPM:
			f.FieldMatchType = &p4.FieldMatch_Lpm{
				Lpm: &p4.FieldMatch_LPM{
					Value: value,
				},
			}
		case p4_config.MatchField_TERNARY:
			f.FieldMatchType = &p4.FieldMatch_Ternary_{
				Ternary: &p4.FieldMatch_Ternary{
					Value: value,
					Mask:  value,
				},
			}
		case p4_config.MatchField_RANGE:
			f.FieldMatchType = &p4.FieldMatch_Range_{
				Range: &p4.FieldMatch_Range{
					Low:  value,
					High: value,
				},
			}
		default:
			//panic
		}
		matchFields[i] = f
	}

	e := &p4.Entity{
		Entity: &p4.Entity_TableEntry{
			TableEntry: &p4.TableEntry{
				TableId:            t.Id,
				Match:              matchFields,
				Action:             action,
				Priority:           priority,
				ControllerMetadata: metadata,
			},
		},
	}
	return e, nil
}

func (t Table) String() string {
	return fmt.Sprintf("Table{Id:%d, Name:%s}", t.Id, t.Name)
}
