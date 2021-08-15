package game

// Role : 각 직업들의 정보를 담고 있는 스트럭처
type Role interface {
	// SendUserSelectGuide 직업 능력을 발휘하기 위한 선택지를 보내는 함수
	SendUserSelectGuide(player *User, g *Game, pageNum int) (msgID string)
	// Action 각 직업별 행동 함수를 다르게 정의하기 위한 함수 선언
	Action(tar *TargetObject, player *User, g *Game)
	// GenLog 함수는 각 직업의 행동 로그를 쌓는 함수이다.
	GenLog(tar *TargetObject, player *User, g *Game)
	// String 함수는 각 직업명을 리턴하는 함수이다.
	String() string
	// ID 함수는 각 직업의 고유값을 리턴하는 함수이다.
	ID() int
}

const (
	sentinel = iota
	doppelganger
	werewolf
	minion
	freemason
	seer
	robber
	troublemaker
	drunk
	insomniac
	hunter
	villager
	tanner
)

// GenerateRole 해당 함수 수정시에 role_guide.json도 수정이 필요
func GenerateRole(num int) (r Role) {
	switch num {
	case sentinel:
		r = &Sentinel{num}
	case doppelganger:
		r = &Doppelganger{num}
	case werewolf:
		r = &Werewolf{num}
	case minion:
		r = &Minion{num}
	case freemason:
		r = &Freemason{num}
	case seer:
		r = &Seer{num}
	case robber:
		r = &Robber{num}
	case troublemaker:
		r = &TroubleMaker{num}
	case drunk:
		r = &Drunk{num}
	case insomniac:
		r = &Insomniac{num}
	case hunter:
		r = &Hunter{num}
	case villager:
		r = &Villager{num}
	case tanner:
		r = &Tanner{num}
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
	//  2:   o     x        x	GetRole, setRole, SetProtect
	//  3:   x     x        o	GetDisRole, setDisRole, GetRoleUsers
	// -1:   x     x        x	RotateAllUserRole, GetRoleUsers
	uid1       string
	uid2       string
	disRoleIdx int
}
