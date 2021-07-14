package game

import (
	"strconv"

	embed "github.com/clinet/discordgo-embed"
)

// Seer 는 한밤의 늑대인간 중 <예언자> 에 대한 객체이다.
type Seer struct {
	id int
}

// SendUserSelectGuide 직업 능력을 발휘하기 위한 선택지를 보내는 함수
func (sr *Seer) SendUserSelectGuide(player *User, g *Game, pageNum int) string {
	title := ""
	if pageNum == 0 {
		title += "직업을 알아낼 플레이어를 고르세요"
	} else {
		title += "세 개의 직업 중 보지 않을 직업을 고르세요"
		curEmbed := embed.NewEmbed()
		curEmbed.SetTitle(title)
		curEmbed.AddField("버려진 직업 셋 중 하나를 선택해 나머지 직업들을 볼 수 있습니다.", "1번 🃏 2번 🃏 3번 🃏")
		msgObj, _ := g.Session.ChannelMessageSendEmbed(player.dmChanID, curEmbed.MessageEmbed)
		for i := 0; i < 3; i++ {
			g.Session.MessageReactionAdd(player.dmChanID, msgObj.ID, g.Emj["n"+strconv.Itoa(i+1)])
		}
		return msgObj.ID
	}
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
	g.Session.MessageReactionAdd(player.dmChanID, msgObj.ID, g.Emj["DISCARD"])
	return msgObj.ID
}

// Action 함수는 <예언자> 의 특수능력 사용에 대한 함수이다.
func (sr *Seer) Action(tar *TargetObject, player *User, g *Game) {
	switch tar.actionType {
	case 2:
		role := g.GetRole(tar.uid1)
		msg := "`" + g.FindUserByUID(tar.uid1).nick + "` 의 직업은 "
		msg += "`" + role.String() + "` 입니다."
		g.Session.ChannelMessageSendEmbed(player.dmChanID, embed.NewGenericEmbed("직업 확인", msg))
	case 3:
		msg := ""
		for i := 0; i < 3; i++ {
			if i != tar.disRoleIdx {
				role := g.GetDisRole(i)
				msg += "`" + role.String() + "` "
			}
		}
		msg += "이(가) 버려져 있습니다."
		g.Session.ChannelMessageSendEmbed(player.dmChanID, embed.NewGenericEmbed("직업 확인", msg))
	}
}

// GenLog 함수는 <예언자> 의 특수능력 사용에 대한 함수이다.
func (sr *Seer) GenLog(tar *TargetObject, player *User, g *Game) {
	msg := ""
	switch tar.actionType {
	case 2:
		role := g.GetRole(tar.uid1)
		msg += "예언자 `" + player.nick + "` 는 "
		msg += "`" + g.FindUserByUID(tar.uid1).nick + "` 의 직업 `" + role.String() + "` 을(를) 확인했습니다."
	case 3:
		msg += "예언자 `" + player.nick + "` 는 버려진 카드 "
		for i := 0; i < 3; i++ {
			if i != tar.disRoleIdx {
				role := g.GetDisRole(i)
				msg += "`" + role.String() + "` "
			}
		}
		msg += "를 확인했습니다."
	}
	g.AppendLog(msg)
}

// String 함수는 <예언자> 문자열을 반환하는 함수이다.
func (sr *Seer) String() string {
	return "예언자"
}

// ID 함수는 <예언자> 의 고유값을 반환하는 함수이다.
func (sr *Seer) ID() int {
	return sr.id
}
