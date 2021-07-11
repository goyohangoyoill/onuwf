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

	// InitState 함수는 스테이트가 시작할 때 필요한 메세지를 생성하고 채널이나 개인DM으로 메세지를 보낸 후
	// 메세지 객체를 스테이트의 멤버로 저장합니다.
	// 이 함수는 이전 스테이트가 끝나는 시점에 호출되어야 합니다.
	InitState()

	// stateFinish 함수는 현재 state가 끝나고 다음 state로 넘어갈 때 호출되는 함수입니다.
	// game의 CurState 변수에 다음 state를 생성해서 할당해준 다음
	// 다음 state의 InitState() 함수를 이 함수 안에서 호출해야 합니다
	stateFinish()

	// filterReaction 함수는 각 스테이트에서 보낸 메세지에 리액션 했는지 거르는 함수이다.
	// 각 스테이트에서 보낸 메세지의 아이디와 리액션이 온 아이디가 동일한지 확인 및
	// 메세지에 리액션 한 것을 지워주어야 한다.
	filterReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd)
}
