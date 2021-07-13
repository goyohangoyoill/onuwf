package game

import (
	embed "github.com/clinet/discordgo-embed"
)

// RoleTemp 는 한밤의 늑대인간 중 <직업명> 에 대한 객체이다.
type Temp struct {
	id int
}

// SendUserSelectGuide 직업 능력을 발휘하기 위한 선택지를 보내는 함수
func (tmp *Temp) SendUserSelectGuide(player *User, g *Game, pageNum int) string {
	tmpEmbed := embed.NewEmbed()
	msg, _ := g.Session.ChannelMessageSendEmbed(player.dmChanID, tmpEmbed.MessageEmbed)
	return msg.ID
}

// Action 함수는 <직업명> 의 특수능력 사용에 대한 함수이다.
func (tmp *Temp) Action(tar *TargetObject, player *User, g *Game) {
	//			<action Type>
	//
	//      uid1  uid2  disRoleIdx
	//  0:   o     o        x	SwapRoleFromUser, CopyRole
	//  1:   o     x        o	SwapRoleFromDiscard
	//  2:   o     x        x	GetRole, setRole, SetProtect
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
func (tmp Temp) GenLog(tar *TargetObject, player *User, g *Game) {
	g.AppendLog("여기에 로그 메시지를 입력하세요")
}

// String 함수는 <직업명> 문자열을 반환하는 함수이다.
func (tmp *Temp) String() string {
	return "직업명"
}

// ID 함수는 <직업명> 의 고유값을 반환하는 함수이다.
func (tmp *Temp) ID() int {
	return tmp.id
}
