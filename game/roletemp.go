package game

import (
	embed "github.com/clinet/discordgo-embed"
)

// RoleTemp 는 한밤의 늑대인간 중 <직업명> 에 대한 객체이다.
type RoleTemp struct {
	Role
}

// Action 함수는 <직업명> 의 특수능력 사용에 대한 함수이다.
func (tmp RoleTemp) Action(tar TargetObject, player *User, g *Game) {
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

// GenLog 함수는 <직업명> 의 특수능력 사용에 대한 함수이다.
func (tmp RoleTemp) GenLog(tar TargetObject, player *User, g *Game) {
	g.AppendLog("여기에 로그 메시지를 입력하세요")
}

// String 함수는 <직업명> 문자열을 반환하는 함수이다.
func (tmp RoleTemp) String() string {
	return "여기에 직업명을 입력하세요"
}
