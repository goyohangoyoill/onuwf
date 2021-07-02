package game

import (
	embed "github.com/clinet/discordgo-embed"
)

// RoleSeer 는 한밤의 늑대인간 중 <예언자> 에 대한 객체이다.
type RoleSeer struct {
	Role
}

// Action 함수는 <예언자> 의 특수능력 사용에 대한 함수이다.
func (sr RoleSeer) Action(tar TargetObject, player *User, g *Game) {
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
func (sr RoleSeer) GenLog(tar TargetObject, player *User, g *Game) {
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
func (sr RoleSeer) String() string {
	return "예언자"
}
