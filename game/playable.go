package game

import (
	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
)

type Playable struct {
	// state 에서 가지고 있는 game
	g *Game

	// factory 에서 쓰이게 될 role index
	roleIndex int

	// 직업추가 확인용 메세지
	roleAddMsg *discordgo.Message

	// 게임입장 확인용 메세지
	enterGameMsg *discordgo.Message
}

// PressNumBtn 사용자가 숫자 이모티콘을 눌렀을 때 Playable에서 하는 동작
func (sPlayable *Playable) PressNumBtn(s *discordgo.Session, r *discordgo.MessageReactionAdd, num int) {

}

// PressDisBtn 사용자가 버려진 카드 이모티콘을 눌렀을 때 Playable에서 하는 동작
func (sPlayable *Playable) PressDisBtn(s *discordgo.Session, r *discordgo.MessageReactionAdd) {

}

// PressYesBtn 사용자가 yes 이모티콘을 눌렀을 때 Playable에서 하는 동작
func (sPlayable *Playable) PressYesBtn(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// 입장 메세지에서 리액션한거라면
	if r.MessageID == sPlayable.enterGameMsg.ID {
		s.ChannelMessageSendEmbed(sPlayable.g.ChanID, embed.NewGenericEmbed("시작가능", ""))
		// 직업추가 메세지에서 리액션한거라면
	} else if r.MessageID == sPlayable.roleAddMsg.ID {

	}
}

// PressNoBtn 사용자가 No 이모티콘을 눌렀을 때 Playable에서 하는 동작
func (sPlayable *Playable) PressNoBtn(s *discordgo.Session, r *discordgo.MessageReactionAdd) {

}

// PressDirBtn 좌 -1, 우 1 사용자가 좌우 방향 이모티콘을 눌렀을 때 Playable에서 하는 동작
func (sPlayable *Playable) PressDirBtn(s *discordgo.Session, r *discordgo.MessageReactionAdd, dir int) {

}

// SendFinish 사용자가 종료 메세지를 보냈을 때 Playable에서 하는 동작
func (sPlayable *Playable) SendFinish(s *discordgo.Session, m *discordgo.MessageCreate) {

}