// +build linux,amd64,go1.15,!cgo

package role

// Role : 각 직업들의 정보를 담고 있는 스트럭처
type Role interface {
	// Action 각 직업별 행동 함수를 다르게 정의하기 위한 함수 선언
	Action(tar TargetObject, player *User, g *Game)
	String() string
}

type RoleFactory struct {
}

func (rf *RoleFactory) generateRole(num int) (r Role) {
	switch num {
	case 1:
		r = &roleSentinel{}
		/*
			case 2:
				r = roleDoppelganger{}
			case 3:
				r = roleWerewolf{}
			case 4:
				r = roleAlphawolf{}
			case 5:
				r = roleMisticwolf{}
			case 6:
				r = roleMinion{}
			case 7:
				r = roleFreemasonry{}
			case 8:
				r = roleSeer{}
			case 9:
				r = roleApprenticeseer{}
			case 10:
				r = roleParanormalinvestigator{}
			case 11:
				r = roleRober{}
			case 12:
				r = roleWitch{}
			case 13:
				r = roleTroublemaker{}
			case 14:
				r = roleVillageidiot{}
			case 15:
				r = roleDrunk{}
			case 16:
				r = roleInsomniac{}
			case 17:
				r = roleRevealer{}
			case 18:
				r = roleTanner{}
			case 19:
				r = roleHunter{}
			case 20:
				r = roleBodygaurd{}
			case 21:
				r = roleVillager{}
			case 22:
				r = roleDreamwolf{}
		*/

	}
	return r
}

// TargetObject 는 각 직업의 특수 능력이
// 적용되어야 하는 대상을 구분하는 기준으로 사용된다.
type TargetObject struct {
	actionType int
	//			<action Type>
	//
	//      uid1  uid2  disRoleIdx
	//  0:   o     o        x	SwapRoleFromUser, CopyRole
	//  1:   o     x        o	SwapRoleFromDiscard
	//  2:   o     x        x	GetRole, setRole, SetPower
	//  3:   x     x        o	GetDisRole, setDisRole, GetRoleUsers
	// -1:   x     x        x	RotateAllUserRole, GetRoleUsers
	uid1       string
	uid2       string
	disRoleIdx int
}
