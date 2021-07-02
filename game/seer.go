package game

import (
	embed "github.com/clinet/discordgo-embed"
)

type Seer struct {
	ID int
}

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
	}
}

func (sr *Seer) String() string {
	return "예언자"
}
