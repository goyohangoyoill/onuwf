package game

import (
	"strconv"

	embed "github.com/clinet/discordgo-embed"
)

// RoleTemp 는 한밤의 늑대인간 중 <늑대인간> 에 대한 객체이다.
type RoleWerewolf struct {
	Role
}

// Action 함수는 <늑대인간> 의 특수능력 사용에 대한 함수이다.
func (wf RoleWerewolf) Action(tar TargetObject, player *User, g *Game) {
	switch tar.actionType {
	case 1:
		recvRole := g.DisRole[tar.disRoleIdx]
		msg := strconv.Itoa(tar.disRoleIdx) + "번째 버려진 카드는\n"
		msg += "`" + recvRole.String() + "` 입니다."
		g.Session.ChannelMessageSendEmbed(player.dmChanID, embed.NewGenericEmbed("버려진 직업 확인", msg))
	case 0:
		wolves := g.GetRoleUsers(wf)
		//wolves = append(wolves, g.getRoleUsers(roleMisticwolf{})...)
		//wolves = append(wolves, g.getRoleUsers(roleAlphawolf{})...)
		//dreams := g.getRoleUsers(roleDreamwolf{})
		var wolflist string
		//var dreamlist string
		for _, item := range wolves {
			wolflist += "`" + item.nick + "` "
		}
		//for _, item := range dreams {
		//	dreamlist += "`" + item.nick + "` "
		//}
		//wolflist += dreamlist
		msg := "당신의 동료 늑대인간은\n"
		msg += wolflist
		msg += "\n ... 입니다."
		//if len(dreams) == 0 {
		//	msg += "\n\n"
		//	msg += dreamlist + "는 잠에 빠져 서로를 확인하지 못하였지만,"
		//	msg += "당신의 동료 늑대인간입니다."
		//}
		g.Session.ChannelMessageSendEmbed(player.dmChanID, embed.NewGenericEmbed("동료 늑대인간 확인", msg))
	}
}

// GenLog 함수는 <늑대인간> 의 특수능력 사용에 대한 함수이다.
func (wf RoleWerewolf) GenLog(tar TargetObject, player *User, g *Game) {
	g.AppendLog("여기에 로그 메시지를 입력하세요")
	switch tar.actionType {
	case 1:
		recvRole := g.DisRole[tar.disRoleIdx]
		logMsg := "유일한 늑대인간 `" + player.nick + "` 은(는)\n"
		logMsg += "버려진 직업 `" + recvRole.String() + "`를 확인했습니다."
		g.AppendLog(logMsg)
	case 0:
		wolves := g.GetRoleUsers(wf)
		//wolves = append(wolves, g.getRoleUsers(roleMisticwolf{})...)
		//wolves = append(wolves, g.getRoleUsers(roleAlphawolf{})...)
		//dreams := g.getRoleUsers(roleDreamwolf{})
		var wolflist string
		//var dreamlist string
		for _, item := range wolves {
			wolflist += "`" + item.nick + "` "
		}
		//for _, item := range dreams {
		//	dreamlist += "`" + item.nick + "` "
		//}
		//wolflist += dreamlist
		logMsg := "늑대인간인 플레이어들\n"
		logMsg += wolflist
		logMsg += "\n는 서로를 확인하였습니다."
		//if len(dreams) == 0 {
		//	dreamLogMsg += "\n\n"
		//	dreamLogMsg += dreamlist + "는 서로를 확인하지 못하였지만,"
		//	dreamLogMsg += "동료 늑대인간들은 확인하였습니다."
		g.AppendLog(logMsg)
		//g.AppendLog(draemLogMsg)
	}
}

// String 함수는 <늑대인간> 문자열을 반환하는 함수이다.
func (wf RoleWerewolf) String() string {
	return "늑대인간"
}
