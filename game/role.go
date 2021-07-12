package game

// Role : 각 직업들의 정보를 담고 있는 스트럭처
type Role interface {
	// SendUserSelectGuide 직업 능력을 발휘하기 위한 선택지를 보내는 함수
	SendUserSelectGuide(player *User, g *Game, pageNum int) (msgID string)
	// Action 각 직업별 행동 함수를 다르게 정의하기 위한 함수 선언
	Action(tar *TargetObject, player *User, g *Game)
	// GenLog 함수는 각 직업의 행동 로그를 쌓는 함수이다.
	GenLog(tar *TargetObject, player *User, g *Game)
	// Stirng 함수는 각 직업명을 리턴하는 함수이다.
	String() string
	// ID 함수는 각 직업의 고유값을 리턴하는 함수이다.
	ID() int
}

// RoleGuide has info of each role
type RoleGuide struct {
	RoleName  string   `json:"roleName"`
	RoleGuide []string `json:"roleGuide"`
	Max       int      `json:"max"`
	Faction   string   `json:"faction"`
	Priority  int      `json:"priority"`
}

func GenerateRole(num int) (r Role) {
	switch num {
	case 0:
		r = &Sentinel{num}
		/*
			case 1:
				r = &Doppelganger{num}
		*/
	case 2:
		r = &Werewolf{num}
		/*
			case 3:
				r = &Alphawolf{num}
			case 4:
				r = &Misticwolf{num}
		*/
	case 3:
		r = &Minion{num}
		/*
			case 6:
				r = &Freemasonry{num}
		*/
	case 5:
		r = &Seer{num}
		/*
			case 8:
				r = &Apprenticeseer{num}
			case 9:
				r = &Paranormalinvestigator{num}
		*/
	case 10:
		r = &Robber{num}
		/*
			case 11:
				r = &Witch{num}
		*/
	case 7:
		r = &TroubleMaker{num}
	case 11:
		r = &Hunter{num}
	case 12:
		r = &Villager{num}
		/*
			case 13:
				r = &Villageidiot{num}
			case 14:
				r = &Drunk{num}
		*/
		//case 15:
		//	r = &Insomniac{num}
		/*
			case 16:
				r = &Revealer{num}
			case 17:
				r = &Tanner{num}
			case 19:
				r = &Bodygaurd{num}
			case 21:
				r = &Dreamwolf{num}
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
