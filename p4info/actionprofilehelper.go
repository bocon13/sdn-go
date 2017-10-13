package p4info

import (
	"fmt"
	"github.com/bocon13/sdn-go/proto/p4"
	"github.com/bocon13/sdn-go/proto/p4/config"
)

type ActionProfile struct {
	info          *P4Info
	Id            uint32
	Name          string
	ActionProfile *p4_config.ActionProfile
}

func (info *P4Info) GetActionProfile(name string) (ActionProfile, error) {
	for i := range info.Info.ActionProfiles {
		ap := info.Info.ActionProfiles[i]
		if p, ok := matches(ap, name); ok {
			return ActionProfile{info, p.Id, name, ap}, nil
		}
	}
	return ActionProfile{}, fmt.Errorf("action profile %s not found", name)
}

func GetEntityFromActionProfileMember(member *p4.ActionProfileMember) *p4.Entity {
	return &p4.Entity{
		Entity: &p4.Entity_ActionProfileMember{ActionProfileMember: member},
	}
}

func (ap *ActionProfile) GetActionProfileMember(id uint32, a *p4.Action) *p4.ActionProfileMember {
	return &p4.ActionProfileMember{
		ActionProfileId: ap.Id,
		MemberId:        id,
		Action:          a,
	}
}

// FIXME wip
//func (ap *ActionProfile) GetActionProfileGroup(groupId uint32, groupType p4.ActionProfileGroup_Type,
//											   members []interface{}, weights map[uint]) *p4.Entity {
//	apms := make([]*p4.ActionProfileGroup_Member, len(members))
//	for i := range members {
//		switch v := members[i].(type) {
//		case *p4.Entity:
//			fmt.Printf("%q is %v bytes long\n", v, len(v))
//		case *p4.ActionProfileMember:
//			fmt.Printf("%q is %v bytes long\n", v, len(v))
//
//		default:
//			fmt.Printf("I don't know about type %T!\n", v)
//		}
//		members[i] = &p4.ActionProfileGroup_Member{
//			MemberId: id,
//			Weight: 1, //FIXME
//		}
//	}
//
//
//
//	return &p4.Entity{
//		Entity: &p4.Entity_ActionProfileGroup{
//			ActionProfileGroup: &p4.ActionProfileGroup{
//				ActionProfileId: ap.Id,
//				GroupId: groupId,
//				Type: groupType,
//				Members: apms,
//				MaxSize: int32(len(apms)),
//			},
//		},
//	}
//}
