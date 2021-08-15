package game

import (
	"strconv"

	embed "github.com/clinet/discordgo-embed"
)

// Sentinel 은 한밤의 늑대인간 중 <수호자>에 대한 객체이다.
type Sentinel struct {
	id int
}

// SendUserSelectGuide 직업 능력을 발휘하기 위한 선택지를 보내는 함수
func (r *Sentinel) SendUserSelectGuide(player *User, g *Game, _ int) (msgID string) {
	curEmbed := embed.NewEmbed()
	curEmbed.SetTitle(g.Emj["SHIELD"] + " 수호할 플레이어를 고르세요")
	for uIdx, user := range g.UserList {
		curEmbed.AddField(strconv.Itoa(uIdx+1)+"번", user.nick)
	}
	curEmbed.InlineAllFields()
	msgObj, _ := g.Session.ChannelMessageSendEmbed(player.dmChanID, curEmbed.MessageEmbed)
	for i := 0; i < len(g.UserList); i++ {
		g.Session.MessageReactionAdd(player.dmChanID, msgObj.ID, g.Emj["n"+strconv.Itoa(i+1)])
	}
	g.Session.MessageReactionAdd(player.dmChanID, msgObj.ID, g.Emj["DISCARD"])
	return msgObj.ID
}

// Action 함수는 <수호자>의 특수능력 사용에 대한 함수이다.
func (r *Sentinel) Action(tar *TargetObject, player *User, g *Game) {
	tmpEmbed := embed.NewEmbed()
	// 항상 tar.actionType == 2
	tarUser := g.FindUserByUID(tar.uid1)
	g.SetProtect(tar.uid1)
	tmpEmbed.SetTitle("빛의 힘이 깃든 방패를 사용하였습니다")
	tmpEmbed.AddField("`"+tarUser.nick+"` 을(를) 수호하였습니다", "누구도 `"+tarUser.nick+"`에게 간섭할 수 없습니다 그가 늑대인간일지라도...")
	g.Session.ChannelMessageSendEmbed(player.dmChanID, tmpEmbed.MessageEmbed)
}

// GenLog 함수는 <수호자>의 특수능력 사용에 대한 함수이다.
func (r *Sentinel) GenLog(tar *TargetObject, player *User, g *Game) {
	var msg string
	// 항상 tar.actionType == 2
	tarUser := g.FindUserByUID(tar.uid1)
	msg = "수호자 `" + player.nick + "` 은(는) "
	if tarUser == nil {
		msg += "아무도 수호하지 않았습니다"
	} else {
		msg += "`" + tarUser.nick + "` 을(를) 방패로 수호하였습니다"
	}
	g.AppendLog(msg)
}

// String 함수는 <수호자> 문자열을 반환하는 함수이다.
func (r *Sentinel) String() string {
	return "수호자"
}

// ID 함수는 <수호자>의 고유값을 반환하는 함수이다.
func (r *Sentinel) ID() int {
	return r.id
}
