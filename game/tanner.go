package game

// Tanner 는 한밤의 늑대인간 중 <무두장이> 에 대한 객체이다.
type Tanner struct {
	id int
}

// SendUserSelectGuide 직업 능력을 발휘하기 위한 선택지를 보내는 함수
func (tmp *Tanner) SendUserSelectGuide(player *User, g *Game, pageNum int) string {
	return "Tanner has no special msg"
}

// Action 함수는 <무두장이> 의 특수능력 사용에 대한 함수이다.
func (tmp *Tanner) Action(tar *TargetObject, player *User, g *Game) {
	//			<action Type>
	//
	//      uid1  uid2  disRoleIdx
	//  0:   o     o        x	SwapRoleFromUser, CopyRole
	//  1:   o     x        o	SwapRoleFromDiscard
	//  2:   o     x        x	GetRole, setRole, SetPower
	//  3:   x     x        o	GetDisRole, setDisRole, GetRoleUsers
	// -1:   x     x        x	RotateAllUserRole, GetRoleUsers
}

// GenLog 함수는 <무두장이> 의 특수능력 사용에 대한 함수이다.
func (tmp Tanner) GenLog(tar *TargetObject, player *User, g *Game) {
}

// String 함수는 <무두장이> 문자열을 반환하는 함수이다.
func (tmp *Tanner) String() string {
	return "무두장이"
}

// ID 함수는 <무두장이> 의 고유값을 반환하는 함수이다.
func (tmp *Tanner) ID() int {
	return tmp.id
}
