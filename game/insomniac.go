package game

import (
	embed "github.com/clinet/discordgo-embed"
)

// RoleTemp 는 한밤의 늑대인간 중 <직업명> 에 대한 객체이다.
type Insomniac struct {
	id int
}

func (is *Insomniac) SendUserSelectGuide(player *User, g *Game, pageNum int) string {
	return "Insomniac have no Msg"
}

// Action 함수는 <직업명> 의 특수능력 사용에 대한 함수이다.
func (is *Insomniac) Action(tar *TargetObject, player *User, g *Game) {
	//			<action Type>
	//
	//      uid1  uid2  disRoleIdx
	//  0:   o     o        x	SwapRoleFromUser, CopyRole
	//  1:   o     x        o	SwapRoleFromDiscard
	//  2:   o     x        x	GetRole, setRole, SetPower
	//  3:   x     x        o	GetDisRole, setDisRole, GetRoleUsers
	// -1:   x     x        x	RotateAllUserRole, GetRoleUsers
	switch tar.actionType {
	case 2:
		role := g.GetRole(tar.uid1)
		msg := "`당신의 직업은 "
		msg += "`" + role.String() + "`입니다."
		g.Session.ChannelMessageSendEmbed(player.dmChanID, embed.NewGenericEmbed("직업 확인", msg))
	}
}

// GenLog 함수는 <직업명> 의 특수능력 사용에 대한 함수이다.
func (is Insomniac) GenLog(tar *TargetObject, player *User, g *Game) {
	msg := ""
	switch tar.actionType {
	case 2:
		role := g.GetRole(tar.uid1)
		msg += "불면증환자 " + player.nick + "` 는 "
		msg += "`자신의 직업 `" + role.String() + "`을 (를) 확인했습니다."
	}
	g.AppendLog(msg)
}

// String 함수는 <직업명> 문자열을 반환하는 함수이다.
func (is *Insomniac) String() string {
	return "불면증환자"
}

// ID 함수는 <직업명> 의 고유값을 반환하는 함수이다.
func (is *Insomniac) ID() int {
	return is.id
}
