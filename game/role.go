package game

// Role : 각 직업들의 정보를 담고 있는 스트럭처
type Role interface {
	// Action 각 직업별 행동 함수를 다르게 정의하기 위한 함수 선언
	Action(tar TargetObject, player *User, g *Game)
	// GenLog 함수는 각 직업의 행동 로그를 쌓는 함수이다.
	GenLog(tar TargetObject, player *User, g *Game)
	// Stirng 함수는 각 직업명을 리턴하는 함수이다.
	String() string
}

// RoleGuide has info of each role
type RoleGuide struct {
	RoleName  string   `json:"roleName"`
	RoleGuide []string `json:"roleGuide"`
	Max       int      `json:"max"`
	Faction   string   `json:"faction"`
}

func GenerateRole(num int) (r Role) {
	switch num {
	case 0:
		r = RoleSentinel{}
		/*
			case 1:
				r = Doppelganger{}
		*/
	case 2:
		r = RoleWerewolf{}
		/*
			case 3:
				r = Alphawolf{}
			case 4:
				r = Misticwolf{}
			case 5:
				r = Minion{}
			case 6:
				r = Freemasonry{}
		*/
	case 7:
		r = RoleSeer{}
		/*
			case 8:
				r = Apprenticeseer{}
			case 9:
				r = Paranormalinvestigator{}
			case 10:
				r = Rober{}
			case 11:
				r = Witch{}
		*/
	case 12:
		r = RoleTroubleMaker{}
		/*
			case 13:
				r = Villageidiot{}
			case 14:
				r = Drunk{}
			case 15:
				r = Insomniac{}
			case 16:
				r = Revealer{}
			case 17:
				r = Tanner{}
			case 18:
				r = Hunter{}
			case 19:
				r = Bodygaurd{}
			case 20:
				r = Villager{}
			case 21:
				r = Dreamwolf{}
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
