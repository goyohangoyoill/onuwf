package game

import (
	embed "github.com/clinet/discordgo-embed"
)

// TroubleMaker 는 한밤의 늑대인간 중 <말썽쟁이> 에 대한 객체이다.
type TroubleMaker struct {
	id int
}

// Action 함수는 <말썽쟁이> 의 특수능력 사용에 대한 함수이다.
func (tm *TroubleMaker) Action(tar *TargetObject, player *User, g *Game) {
	//			<action Type>
	//
	//      uid1  uid2  disRoleIdx
	//  0:   o     o        x	SwapRoleFromUser, CopyRole
	//  1:   o     x        o	SwapRoleFromDiscard
	//  2:   o     x        x	GetRole, setRole, SetPower
	//  3:   x     x        o	GetDisRole, setDisRole, GetRoleUsers
	// -1:   x     x        x	RotateAllUserRole, GetRoleUsers
	switch tar.actionType {
	case 0:
		g.SwapRoleFromUser(tar.uid1, tar.uid2)
		user1 := g.FindUserByUID(tar.uid1)
		user2 := g.FindUserByUID(tar.uid2)
		msg := "`" + user1.nick + "`, `" + user2.nick + "`"
		msg += " 의 직업을 맞바꿨습니다."
		g.Session.ChannelMessageSendEmbed(player.dmChanID, embed.NewGenericEmbed("능력 사용", msg))
	}
}

// GenLog 함수는 <말썽쟁이> 의 특수능력 사용에 대한 함수이다.
func (tm *TroubleMaker) GenLog(tar *TargetObject, player *User, g *Game) {
	switch tar.actionType {
	case 0:
		user1 := g.FindUserByUID(tar.uid1)
		user2 := g.FindUserByUID(tar.uid2)
		role1 := g.GetRole(tar.uid1)
		role2 := g.GetRole(tar.uid2)

		msg := "말썽쟁이 `" + player.nick + "` 는,\n"
		msg += "(`" + role1.String() + "`) `" + user1.nick + "`, (`" + role2.String() + "`) `" + user2.nick + "`\n"
		msg += "의 직업을 맞바꿨습니다."
		g.AppendLog(msg)
	}
}

// String 함수는 <말썽쟁이> 문자열을 반환하는 함수이다.
func (tm *TroubleMaker) String() string {
	return "말썽쟁이"
}

// ID 함수는 <말썽쟁이> 의 고유값을 반환하는 함수이다.
func (tm *TroubleMaker) ID() int {
	return tm.id
}
