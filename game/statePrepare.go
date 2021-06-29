package game

import "github.com/bwmarrin/discordgo"

// Prepare is test
type Prepare struct {
	g            *Game
	roleIndex    int
	roleAddMsg   *discordgo.Message
	enterGameMsg *discordgo.Message
}

// PressNumBtn 사용자가 숫자 이모티콘을 눌렀을 때 state에서 하는 동작
func (p Prepare) PressNumBtn(s *discordgo.Session, r *discordgo.MessageReactionAdd, num int) {
}

// PressDisBtn 사용자가 버려진 카드 이모티콘을 눌렀을 때 state에서 하는 동작
func (p Prepare) PressDisBtn(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
}

// PressYesBtn 사용자가 yes 이모티콘을 눌렀을 때 state에서 하는 동작
func (p Prepare) PressYesBtn(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
}

// PressNoBtn 사용자가 No 이모티콘을 눌렀을 때 state에서 하는 동작
func (p Prepare) PressNoBtn(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
}

// PressDirBtn 좌 -1, 우 1 사용자가 좌우 방향 이모티콘을 눌렀을 때 state에서 하는 동작
func (p Prepare) PressDirBtn(s *discordgo.Session, r *discordgo.MessageReactionAdd, dir int) {
}
