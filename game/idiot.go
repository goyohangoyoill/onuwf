package game

import (
	embed "github.com/clinet/discordgo-embed"
)

// RoleIdiot 는 한밤의 늑대인간 중 <동네바보> 에 대한 객체이다.
type Idiot struct {
	id int
}

// SendUserSelectGuide 직업 능력을 발휘하기 위한 선택지를 보내는 함수
func (idt *Idiot) SendUserSelectGuide(player *User, g *Game, pageNum int) string {
	idtEmbed := embed.NewEmbed()
	footer := g.Emj["YES"] + " : 능력 사용하기\n" + g.Emj["NO"] + " : 능력 사용하지 않기"
	switch pageNum {
	case 0:
		idtEmbed.SetTitle("능력 사용")
		idtEmbed.AddField("모든 플레이어의 직업을 입장순서대로\n다음 사람에게 넘깁니다.", footer)
	case 1:
		idtEmbed.SetTitle("진짜로요?")
		idtEmbed.AddField("모든 플레이어의 직업을 입장순서대로\n다음 사람에게 넘깁니다.", footer)
	case 2:
		idtEmbed.SetTitle("바보짓 장전!")
		idtEmbed.AddField("마지막으로 묻습니다 후회 안하실건가요?", footer)
	}
	msg, _ := g.Session.ChannelMessageSendEmbed(player.dmChanID, idtEmbed.MessageEmbed)
	g.Session.MessageReactionAdd(player.dmChanID, msg.ID, g.Emj["YES"])
	g.Session.MessageReactionAdd(player.dmChanID, msg.ID, g.Emj["NO"])
	return msg.ID
}

// Action 함수는 <동네바보> 의 특수능력 사용에 대한 함수이다.
func (idt *Idiot) Action(tar *TargetObject, player *User, g *Game) {
	//			<action Type>
	//
	//      uid1  uid2  disRoleIdx
	//  0:   o     o        x	SwapRoleFromUser, CopyRole
	//  1:   o     x        o	SwapRoleFromDiscard
	//  2:   o     x        x	GetRole, setRole, SetProtect
	//  3:   x     x        o	GetDisRole, setDisRole, GetRoleUsers
	// -1:   x     x        x	RotateAllUserRole, GetRoleUsers
	idtEmbed := embed.NewGenericEmbed("hello", "bye")
	switch tar.actionType {
	case -1:
		g.RotateAllUserRole()
	}
	g.Session.ChannelMessageSendEmbed(player.dmChanID, idtEmbed)
}

// GenLog 함수는 <동네바보> 의 특수능력 사용에 대한 함수이다.
func (idt Idiot) GenLog(tar *TargetObject, player *User, g *Game) {
	g.AppendLog("동네 바보`" + player.nick + "`(이)가 일을 저질렀습니다. 헤헤")
	// RotateAllUserRole 에서 로그도 쌓는다.
}

// String 함수는 <동네바보> 문자열을 반환하는 함수이다.
func (idt *Idiot) String() string {
	return "동네바보"
}

// ID 함수는 <동네바보> 의 고유값을 반환하는 함수이다.
func (idt *Idiot) ID() int {
	return idt.id
}
