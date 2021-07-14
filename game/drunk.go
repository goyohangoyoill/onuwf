package game

import (
	"strconv"

	embed "github.com/clinet/discordgo-embed"
)

// RoleTemp 는 한밤의 늑대인간 중 <직업명> 에 대한 객체이다.
type Drunk struct {
	id int
}

// SendUserSelectGuide 직업 능력을 발휘하기 위한 선택지를 보내는 함수
func (dr *Drunk) SendUserSelectGuide(player *User, g *Game, pageNum int) string {
	curEmbed := embed.NewEmbed()
	curEmbed.SetTitle("버려진 직업 셋 중 하나를 선택하세요. 선택한 직업과 변경됩니다")
	msgObj, _ := g.Session.ChannelMessageSendEmbed(player.dmChanID, curEmbed.MessageEmbed)
	for i := 0; i < 3; i++ {
		g.Session.MessageReactionAdd(player.dmChanID, msgObj.ID, g.Emj["n"+strconv.Itoa(i+1)])
	}
	return msgObj.ID
}

// Action 함수는 <직업명> 의 특수능력 사용에 대한 함수이다.
func (dr *Drunk) Action(tar *TargetObject, player *User, g *Game) {
	//			<action Type>
	//
	//      uid1  uid2  disRoleIdx
	//  0:   o     o        x	SwapRoleFromUser, CopyRole
	//  1:   o     x        o	SwapRoleFromDiscard
	//  2:   o     x        x	GetRole, setRole, SetPower
	//  3:   x     x        o	GetDisRole, setDisRole, GetRoleUsers
	// -1:   x     x        x	RotateAllUserRole, GetRoleUsers

	switch tar.actionType {
	case 1:
		// do smthing
		g.SwapRoleFromDiscard(player.UserID, tar.disRoleIdx)

	}

}

// GenLog 함수는 <직업명> 의 특수능력 사용에 대한 함수이다.
func (dr *Drunk) GenLog(tar *TargetObject, player *User, g *Game) {
	msg := ""
	switch tar.actionType {
	case 1:
		Orirole := g.GetOriRole(tar.uid1)
		role := g.GetRole(tar.uid1)
		msg += "주정뱅이 `" + player.nick + "` 는 "
		msg += "자신의 직업 `" + Orirole.String() + "`을(를) 버려진 카드`" + role.String() + "`와(과) 교환했습니다."
	}
	g.AppendLog(msg)

}

// String 함수는 <직업명> 문자열을 반환하는 함수이다.
func (dr *Drunk) String() string {
	return "주정뱅이"
}

// ID 함수는 <직업명> 의 고유값을 반환하는 함수이다.
func (dr *Drunk) ID() int {
	return dr.id
}
