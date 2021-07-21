package game

import (
	"strconv"

	embed "github.com/clinet/discordgo-embed"
)

// Robber 는 한밤의 늑대인간 중 <강도> 에 대한 객체이다.
type Robber struct {
	id int
}

// SendUserSelectGuide 직업 능력을 발휘하기 위한 선택지를 보내는 함수
func (rb *Robber) SendUserSelectGuide(player *User, g *Game, pageNum int) string {
	title := "직업을 훔칠 플레이어를 고르세요"
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

// Action 함수는 <강도> 의 특수능력 사용에 대한 함수이다.
func (rb *Robber) Action(tar *TargetObject, player *User, g *Game) {
	//			<action Type>
	//
	//      uid1  uid2  disRoleIdx
	//  0:   o     o        x	SwapRoleFromUser, CopyRole
	//  1:   o     x        o	SwapRoleFromDiscard
	//  2:   o     x        x	GetRole, setRole, SetProtect
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
func (rb *Robber) GenLog(tar *TargetObject, player *User, g *Game) {
	var msg string
	switch tar.actionType {
	case 2:
		tarRole := g.GetRole(tar.uid1)
		tarUser := g.FindUserByUID(tar.uid1)
		msg = "강도 `" + player.nick + "` 은(는) `" + tarUser.nick + "` 의 직업 `" + tarRole.String() + "` 을(를) 확인하고 자신의 직업과 맞바꾸었습니다."
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
