package game

import (
	"strconv"

	embed "github.com/clinet/discordgo-embed"
)

// RoleFreemason 는 한밤의 늑대인간 중 <프리메이슨> 에 대한 객체이다.
type Freemason struct {
	id int
}

// SendUserSelectGuide 직업 능력을 발휘하기 위한 선택지를 보내는 함수
func (frm *Freemason) SendUserSelectGuide(player *User, g *Game, pageNum int) string {
	return "freemason has no guide msg"
}

// Action 함수는 <프리메이슨> 의 특수능력 사용에 대한 함수이다.
func (frm *Freemason) Action(tar *TargetObject, player *User, g *Game) {
	//			<action Type>
	//
	//      uid1  uid2  disRoleIdx
	//  0:   o     o        x	SwapRoleFromUser, CopyRole
	//  1:   o     x        o	SwapRoleFromDiscard
	//  2:   o     x        x	GetRole, setRole, SetProtect
	//  3:   x     x        o	GetDisRole, setDisRole, GetRoleUsers
	// -1:   x     x        x	RotateAllUserRole, GetRoleUsers
	frmEmbed := embed.NewEmbed()
	frmEmbed.SetTitle("동료 프리메이슨 확인")
	switch tar.actionType {
	case -1:
		frList := g.GetOriRoleUsers(frm)
		if len(frList) == 1 {
			frmEmbed.AddField("당신은 유일한 프리메이슨입니다.", "당신은 동료 프리메이슨을 확인하지 못했습니다.")
		} else {
			for i, user := range frList {
				frmEmbed.AddField(strconv.Itoa(i+1)+"번", "`"+user.nick+"`")
			}
			frmEmbed.InlineAllFields()
		}
	}
	g.Session.ChannelMessageSendEmbed(player.dmChanID, frmEmbed.MessageEmbed)
}

// GenLog 함수는 <프리메이슨> 의 특수능력 사용에 대한 함수이다.
func (frm *Freemason) GenLog(tar *TargetObject, player *User, g *Game) {
	var frmMsg string
	switch tar.actionType {
	case -1:
		frList := g.GetOriRoleUsers(frm)
		if len(frList) == 1 {
			frmMsg = "프리메이슨 `" + frList[0].nick + "`은(는)"
			frmMsg += "자신이 유일한 프리메이슨임을 확인했습니다."
		} else {
			frmMsg = "프리메이슨 "
			for _, user := range frList {
				frmMsg += "`" + user.nick + "` "
			}
			frmMsg += "은(는) 서로를 확인했습니다."
		}
	}
	g.AppendLog(frmMsg)
}

// String 함수는 <프리메이슨> 문자열을 반환하는 함수이다.
func (frm *Freemason) String() string {
	return "프리메이슨"
}

// ID 함수는 <프리메이슨> 의 고유값을 반환하는 함수이다.
func (frm *Freemason) ID() int {
	return frm.id
}
