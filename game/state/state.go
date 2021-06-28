// +build linux,amd64,go1.15,!cgo

package state

import (
	"github.com/bwmarrin/discordgo"
)

type State interface {
	// 사용자 인원수 3 ~ 26
	// num: 0 ~ 23
	// PressNumBtn 사용자가 숫자 이모티콘을 눌렀을 때 state에서 하는 동작
	PressNumBtn(s *discordgo.Session, r *discordgo.MessageReactionAdd, num int)

	// PressDisBtn 사용자가 버려진 카드 이모티콘을 눌렀을 때 state에서 하는 동작
	PressDisBtn(s *discordgo.Session, r *discordgo.MessageReactionAdd)

	// PressYesBtn 사용자가 yes 이모티콘을 눌렀을 때 state에서 하는 동작
	PressYesBtn(s *discordgo.Session, r *discordgo.MessageReactionAdd)

	// PressNoBtn 사용자가 No 이모티콘을 눌렀을 때 state에서 하는 동작
	PressNoBtn(s *discordgo.Session, r *discordgo.MessageReactionAdd)

	// PressDirBtn 좌 -1, 우 1 사용자가 좌우 방향 이모티콘을 눌렀을 때 state에서 하는 동작
	PressDirBtn(s *discordgo.Session, r *discordgo.MessageReactionAdd, dir int)
}
