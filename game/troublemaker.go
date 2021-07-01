package game

import (
	embed "github.com/clinet/discordgo-embed"
)

type RoleTroubleMaker struct {
	Role
}

func (tm *RoleTroubleMaker) Action(tar TargetObject, player *User, g *Game) {
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
		g.Session.ChannelMessageSendEmbed(player.dmChanID, embed.NewGenericEmbed("능력 사용", msg))
	}
}

func (tm *RoleTroubleMaker) String() string {
	return "말썽쟁이"
}
