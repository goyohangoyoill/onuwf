package game

import (
	"github.com/bwmarrin/discordgo"
)

// State 는 리액션 입력이 발생했을 때 현재 상태에 따라 다른 함수를 호출하는 작업을 수행하는
// 인터페이스로, 숫자 버튼, 쓰레기통 버튼, 예/아니오 버튼, 방향 버튼을 인식할 수 있게 구현한다.
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
