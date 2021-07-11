package game

import embed "github.com/clinet/discordgo-embed"

// Sentinel은 한밤의 늑대인간 중 <수호자>에 대한 객체이다.
type Sentinel struct {
	id int
}

// SendUserSelectGuide 직업 능력을 발휘하기 위한 선택지를 보내는 함수
func (r *Sentinel) SendUserSelectGuide(player *User, g *Game, pageNum int) (msgID string) {
	return ""
}

// Action 함수는 <수호자>의 특수능력 사용에 대한 함수이다.
func (r *Sentinel) Action(tar *TargetObject, player *User, g *Game) {
	//session 메세지는 state에서 보낼거임
	//action에서는 game 상태 바꾸는 action만

	//			<action Type>
	//
	//      uid1  uid2  disRoleIdx
	//  0:   o     o        x	SwapRoleFromUser, CopyRole
	//  1:   o     x        o	SwapRoleFromDiscard
	//  2:   o     x        x	GetRole, setRole, SetPower
	//  3:   x     x        o	GetDisRole, setDisRole, GetRoleUsers
	// -1:   x     x        x	RotateAllUserRole, GetRoleUsers
	tmpEmbed := embed.NewGenericEmbed("hello", "bye")
	switch tar.actionType {
	case 1:
		// do smthing
	case 2:
		// do smthing
	}
	g.Session.ChannelMessageSendEmbed("temp", tmpEmbed)
}

// GenLog 함수는 <수호자>의 특수능력 사용에 대한 함수이다.
func (r *Sentinel) GenLog(tar *TargetObject, player *User, g *Game) {
	g.AppendLog("여기에 로그 메시지를 입력하세요")
}

// String 함수는 <수호자> 문자열을 반환하는 함수이다.
func (r *Sentinel) String() string {
	return "수호자"
}

// ID 함수는 <수호자>의 고유값을 반환하는 함수이다.
func (r *Sentinel) ID() int {
	return r.id
}
