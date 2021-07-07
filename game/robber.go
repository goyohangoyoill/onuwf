package game

import (
	embed "github.com/clinet/discordgo-embed"
)

// Robber 는 한밤의 늑대인간 중 <강도> 에 대한 객체이다.
type Robber struct {
	id int
}

// Action 함수는 <강도> 의 특수능력 사용에 대한 함수이다.
func (rb *Robber) Action(tar *TargetObject, player *User, g *Game) {
	//			<action Type>
	//
	//      uid1  uid2  disRoleIdx
	//  0:   o     o        x	SwapRoleFromUser, CopyRole
	//  1:   o     x        o	SwapRoleFromDiscard
	//  2:   o     x        x	GetRole, setRole, SetPower
	//  3:   x     x        o	GetDisRole, setDisRole, GetRoleUsers
	// -1:   x     x        x	RotateAllUserRole, GetRoleUsers
	tmpEmbed := embed.NewEmbed()
	switch tar.actionType {
	case 2:
		tarRole := g.GetRole(tar.uid1)
		tarUser := g.FindUserByUID(tar.uid1)
		g.SwapRoleFromUser(tar.uid1, player.UserID)
		tmpEmbed.SetTitle("직업을 훔쳤습니다")
		tmpEmbed.AddField("`"+tarUser.nick+"`의 직업", "`"+tarUser.nick+"`의 직업은 `"+tarRole.String()+"` 였습니다. 하지만 이젠 아니죠...")
	}
	g.Session.ChannelMessageSendEmbed(player.dmChanID, tmpEmbed.MessageEmbed)
}

// GenLog 함수는 <강도> 의 특수능력 사용에 대한 함수이다.
func (rb *Robber) GenLog(tar TargetObject, player *User, g *Game) {
	var msg string
	switch tar.actionType {
	case 2:
		tarRole := g.GetRole(tar.uid1)
		tarUser := g.FindUserByUID(tar.uid1)
		msg = "강도 `" + player.nick + "` 은(는) `" + tarUser.nick + "` 의 직업 `" + tarRole.String() + "` 을(를) 확인하고\n자신의 직업과 맞바꾸었습니다."
	}
	g.AppendLog(msg)
}

// String 함수는 <강도> 문자열을 반환하는 함수이다.
func (rb *Robber) String() string {
	return "강도"
}

// ID 함수는 <강도> 의 고유값을 반환하는 함수이다.
func (rb *Robber) ID() int {
	return rb.id
}
