package game

import (
	"strconv"

	embed "github.com/clinet/discordgo-embed"
)

// Doppelganger 는 한밤의 늑대인간 중 <도플갱어> 에 대한 객체이다.
type Doppelganger struct {
	id int
}

// SendUserSelectGuide 직업 능력을 발휘하기 위한 선택지를 보내는 함수
func (dpl *Doppelganger) SendUserSelectGuide(player *User, g *Game, pageNum int) string {
	title := "직업을 복사할 플레이어를 고르세요"
	curEmbed := embed.NewEmbed()
	curEmbed.SetTitle(title)
	for uIdx, user := range g.UserList {
		if !g.IsProtected(user.UserID) {
			curEmbed.AddField(strconv.Itoa(uIdx+1)+"번", user.nick)
		} else {
			curEmbed.AddField(strconv.Itoa(uIdx+1)+"번", "🛡"+user.nick)
		}
	}
	curEmbed.InlineAllFields()
	msgObj, _ := g.Session.ChannelMessageSendEmbed(player.dmChanID, curEmbed.MessageEmbed)
	for i := 0; i < len(g.UserList); i++ {
		g.Session.MessageReactionAdd(player.dmChanID, msgObj.ID, g.Emj["n"+strconv.Itoa(i+1)])
	}
	return msgObj.ID
}

// Action 함수는 <도플갱어> 의 특수능력 사용에 대한 함수이다.
func (dpl *Doppelganger) Action(tar *TargetObject, player *User, g *Game) {
	//			<action Type>
	//
	//      uid1  uid2  disRoleIdx
	//  0:   o     o        x	SwapRoleFromUser, CopyRole
	//  1:   o     x        o	SwapRoleFromDiscard
	//  2:   o     x        x	GetRole, setRole, SetProtect
	//  3:   x     x        o	GetDisRole, setDisRole, GetRoleUsers
	// -1:   x     x        x	RotateAllUserRole, GetRoleUsers
	dplEmbed := embed.NewGenericEmbed("hello", "bye")
	switch tar.actionType {
	case 1:
		// do smthing
	case 2:
		// do smthing
	}
	g.Session.ChannelMessageSendEmbed("Doppelganger", dplEmbed)
}

// GenLog 함수는 <도플갱어> 의 특수능력 사용에 대한 함수이다.
func (dpl *Doppelganger) GenLog(tar *TargetObject, player *User, g *Game) {
	g.AppendLog("여기에 로그 메시지를 입력하세요")
}

// String 함수는 <도플갱어> 문자열을 반환하는 함수이다.
func (dpl *Doppelganger) String() string {
	return "도플갱어"
}

// ID 함수는 <도플갱어> 의 고유값을 반환하는 함수이다.
func (dpl *Doppelganger) ID() int {
	return dpl.id
}
