package game

import (
	embed "github.com/clinet/discordgo-embed"
)

type RoleTemp struct {
	Role
}

func (tmp *RoleTemp) Action(tar TargetObject, player *User, g *Game) {
	//			<action Type>
	//
	//      uid1  uid2  disRoleIdx
	//  0:   o     o        x	SwapRoleFromUser, CopyRole
	//  1:   o     x        o	SwapRoleFromDiscard
	//  2:   o     x        x	GetRole, setRole, SetPower
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

func (tmp *RoleTemp) String() string {
	return "말썽쟁이"
}
