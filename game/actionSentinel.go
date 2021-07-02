package game

import (
	"github.com/bwmarrin/discordgo"
)

type ActionSentinel struct {
	// state 에서 가지고 있는 game
	g *Game

	// sentinel 능력사용 메세지
	sentinelMsg *discordgo.Message
}

// PressNumBtn 사용자가 숫자 이모티콘을 눌렀을 때 ActionSentinel에서 하는 동작
func (sActionSentinel *ActionSentinel) PressNumBtn(s *discordgo.Session, r *discordgo.MessageReactionAdd, num int) {

}

// PressDisBtn 사용자가 버려진 카드 이모티콘을 눌렀을 때 ActionSentinel에서 하는 동작
func (sActionSentinel *ActionSentinel) PressDisBtn(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// do nothing
}

// PressYesBtn 사용자가 yes 이모티콘을 눌렀을 때 ActionSentinel에서 하는 동작
func (sActionSentinel *ActionSentinel) PressYesBtn(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// do nothing
}

// PressNoBtn 사용자가 No 이모티콘을 눌렀을 때 ActionSentinel에서 하는 동작
func (sActionSentinel *ActionSentinel) PressNoBtn(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// do nothing
}

// PressDirBtn 좌 -1, 우 1 사용자가 좌우 방향 이모티콘을 눌렀을 때 ActionSentinel에서 하는 동작
func (sActionSentinel *ActionSentinel) PressDirBtn(s *discordgo.Session, r *discordgo.MessageReactionAdd, dir int) {
	// do nothing
}

// InitState 함수는 ActionSentinel state가 시작할 때 진짜로 게임이 시작되므로
// game에 UserList에 직업을 랜덤 할당해주고 각 유저에게 직업소개 개인 DM을 보낸 후
// 센티넬 직업을 가진 유저에게 센티넬 동작 DM을 보내고 이를 ActionSentinel 멤버 변수로 저장합니다.
func (sActionSentinel *ActionSentinel) InitState() {

}

// stateFinish 함수는 sentinel role을 가진 user가 능력사용을 끝내고
// ActionSentinel에서 state가 종료되는 시점에서 호출 됩니다.
// 다음 state인 ActionDoppelganger의 InitState() 함수를 호출합니다.
func (sActionSentinel *ActionSentinel) stateFinish(s *discordgo.Session, r *discordgo.MessageReactionAdd) {

}

// filterReaction 함수는 각 스테이트에서 보낸 메세지에 리액션 했는지 거르는 함수이다.
// 각 스테이트에서 보낸 메세지의 아이디와 리액션이 온 아이디가 동일한지 확인 및
// 메세지에 리액션 한 것을 지워주어야 한다.
func (sActionSentinel *ActionSentinel) filterReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd) {

}
