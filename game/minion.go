package game

import (
	"strconv"

	embed "github.com/clinet/discordgo-embed"
)

// Minion 는 한밤의 늑대인간 중 <하수인> 에 대한 객체이다.
type Minion struct {
	id int
}

// SendUserSelectGuide 직업 능력을 발휘하기 위한 선택지를 보내는 함수
func (min *Minion) SendUserSelectGuide(player *User, g *Game, pageNum int) string {
	return "Minion have no Msg"
}

// Action 함수는 <하수인> 의 특수능력 사용에 대한 함수이다.
func (min *Minion) Action(tar *TargetObject, player *User, g *Game) {
	//			<action Type>
	//
	//      uid1  uid2  disRoleIdx
	//  0:   o     o        x	SwapRoleFromUser, CopyRole
	//  1:   o     x        o	SwapRoleFromDiscard
	//  2:   o     x        x	GetRole, setRole, SetProtect
	//  3:   x     x        o	GetDisRole, setDisRole, GetRoleUsers
	// -1:   x     x        x	RotateAllUserRole, GetRoleUsers
	minEmbed := embed.NewEmbed()
	minEmbed.SetTitle("당신의 동료 늑대인간은...")
	switch tar.actionType {
	case -1:
		wfList := g.GetRoleUsers(&Werewolf{})
		// wfList = append(wfList, g.GetRoleUsers(&AlphaWolf{})
		// wfList = append(wfList, g.GetRoleUsers(&MisticWolf{})
		// wfList = append(wfList, g.GetRoleUsers(&DreamWolf{})
		for i, wfUser := range wfList {
			minEmbed.AddField(strconv.Itoa(i+1)+"번", "`"+wfUser.nick+"`")
		}
		if len(wfList) == 0 {
			minEmbed.AddField("없습니다.", "당신이 유일한 늑대인간 진영입니다. 살아남으세요.")
		}
		minEmbed.InlineAllFields()
	}
	g.Session.ChannelMessageSendEmbed(player.dmChanID, minEmbed.MessageEmbed)
}

// GenLog 함수는 <하수인> 의 특수능력 사용에 대한 함수이다.
func (min Minion) GenLog(tar *TargetObject, player *User, g *Game) {
	var msg string
	switch tar.actionType {
	case -1:
		wfList := g.GetRoleUsers(&Werewolf{})
		// wfList = append(wfList, g.GetRoleUsers(&AlphaWolf{})
		// wfList = append(wfList, g.GetRoleUsers(&MisticWolf{})
		// wfList = append(wfList, g.GetRoleUsers(&DreamWolf{})
		msg = "하수인 `" + player.nick + "` 은(는) 동료 늑대인간"
		for _, wfUser := range wfList {
			msg += "`" + wfUser.nick + "` "
		}
		msg += "을(를) 확인하였습니다."
		if len(wfList) == 0 {
			msg = "하수인 `" + player.nick + "` 은(는) 늑대인간이 없음을 확인했습니다."
		}
	}
	g.AppendLog(msg)
}

// String 함수는 <하수인> 문자열을 반환하는 함수이다.
func (min *Minion) String() string {
	return "하수인"
}

// ID 함수는 <하수인> 의 고유값을 반환하는 함수이다.
func (min *Minion) ID() int {
	return min.id
}
