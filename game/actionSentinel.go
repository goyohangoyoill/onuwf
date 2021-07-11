package game

import (
	"github.com/bwmarrin/discordgo"
)

type ActionSentinel struct {
	// state 에서 가지고 있는 game
	g *Game

	// sentinel role을 가진 user들에게 보낸 능력사용 메세지ID
	sentinelMsgsID string

	// 센티넬의 선택을 기다릴 채널
	Choice chan int
}

// PressNumBtn 사용자가 숫자 이모티콘을 눌렀을 때 ActionSentinel에서 하는 동작
func (sActionSentinel *ActionSentinel) PressNumBtn(s *discordgo.Session, r *discordgo.MessageReactionAdd, num int) {
	sActionSentinel.filterReaction(s, r)
	if sActionSentinel.g.UserList[num-1].UserID == r.UserID {
		s.ChannelMessageSend(r.ChannelID, "자기 자신을 선택할 수 없습니다.")
		return
	}
	sActionSentinel.Choice <- num - 1
}

// PressDisBtn 사용자가 버려진 카드 이모티콘을 눌렀을 때 ActionSentinel에서 하는 동작
func (sActionSentinel *ActionSentinel) PressDisBtn(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	sActionSentinel.filterReaction(s, r)
	s.ChannelMessageSend(r.ChannelID, "아무도 방패로 보호하지 않았습니다.")
	sActionSentinel.Choice <- -1
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

// InitState 함수는 ActionSentinel state가 시작할 때
// 센티넬 직업을 가진 유저에게 센티넬 동작 DM을 보내고 이를 ActionSentinel 멤버 변수로 저장합니다.
func (sActionSentinel *ActionSentinel) InitState() {
	// sentinel role을 가지고 있는 유저들에게 능력사용 메세지 보낸 후 MessageID 저장
	g := sActionSentinel.g
	role := &Sentinel{}
	sentinel := g.GetRoleUsers(role)[0]
	sActionSentinel.Choice = make(chan int)
	sActionSentinel.sentinelMsgsID = role.SendUserSelectGuide(sentinel, g, 0)
	input := <-sActionSentinel.Choice
	if input == -1 {
		tar := &TargetObject{2, "", "", -1}
		role.GenLog(tar, sentinel, g)
	} else {
		tar := &TargetObject{2, g.UserList[input-1].UserID, "", -1}
		role.Action(tar, sentinel, g)
		role.GenLog(tar, sentinel, g)
	}
}

// stateFinish 함수는 sentinel role을 가진 user가 능력사용을 끝내고
// ActionSentinel에서 state가 종료되는 시점에서 호출 됩니다.
// 다음 state인 ActionDoppelganger의 InitState() 함수를 호출합니다.
func (sActionSentinel *ActionSentinel) stateFinish() {
	// ActionDoppelganger의 Msg 변수 초기화
	sActionSentinel.g.CurState = NewActionInGameGroup(sActionSentinel.g)
	sActionSentinel.g.CurState.InitState()
}

// filterReaction 함수는 각 스테이트에서 보낸 메세지에 리액션 했는지 거르는 함수이다.
// 각 스테이트에서 보낸 메세지의 아이디와 리액션이 온 아이디가 동일한지 확인 및
// 메세지에 리액션 한 것을 지워주어야 한다.
func (sActionSentinel *ActionSentinel) filterReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// 현재 스테이트에서 보낸 메세지에 리액션한 게 아니면 거름
	if !(r.MessageID == sActionSentinel.sentinelMsgsID) {
		return
	}
	// 메세지에 리약션한 거 지워줌
	s.MessageReactionRemove(r.ChannelID, r.MessageID, r.Emoji.Name, r.UserID)
}
