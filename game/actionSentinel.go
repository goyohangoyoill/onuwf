package game

import (
	"github.com/bwmarrin/discordgo"
)

type ActionSentinel struct {
	// state 에서 가지고 있는 game
	g *Game
	// sentinel role을 가진 user들에게 보낸 능력사용 메세지ID
	sentinelMsgsID map[string]string
	// 센티넬의 선택을 기다릴 채널
	UserChoice chan Choice
}

type Choice struct {
	num  int
	user *User
}

// PressNumBtn 사용자가 숫자 이모티콘을 눌렀을 때 ActionSentinel에서 하는 동작
func (sActionSentinel *ActionSentinel) PressNumBtn(s *discordgo.Session, r *discordgo.MessageReaction, num int) {
	// 게임 진행과 관련된 메세지에 달린 리액션 지운다
	if sActionSentinel.filterReaction(s, r) {
		return
	}
	if sActionSentinel.g.UserList[num-1].UserID == r.UserID {
		s.ChannelMessageSend(r.ChannelID, "자기 자신을 선택할 수 없습니다.")
		return
	}
	sActionSentinel.UserChoice <- Choice{num, sActionSentinel.g.FindUserByUID(r.UserID)}
	s.ChannelMessageDelete(r.ChannelID, sActionSentinel.sentinelMsgsID[r.UserID])
}

// PressDisBtn 사용자가 버려진 카드 이모티콘을 눌렀을 때 ActionSentinel에서 하는 동작
func (sActionSentinel *ActionSentinel) PressDisBtn(s *discordgo.Session, r *discordgo.MessageReaction) {
	// 게임 진행과 관련된 메세지에 달린 리액션 지운다
	if sActionSentinel.filterReaction(s, r) {
		return
	}
	s.ChannelMessageSend(r.ChannelID, "아무도 방패로 보호하지 않았습니다.")
	sActionSentinel.UserChoice <- Choice{-1, sActionSentinel.g.FindUserByUID(r.UserID)}
}

// PressYesBtn 사용자가 yes 이모티콘을 눌렀을 때 ActionSentinel에서 하는 동작
func (sActionSentinel *ActionSentinel) PressYesBtn(s *discordgo.Session, r *discordgo.MessageReaction) {
	// 게임 진행과 관련된 메세지에 달린 리액션 지운다
	if sActionSentinel.filterReaction(s, r) {
		return
	}
}

// PressNoBtn 사용자가 No 이모티콘을 눌렀을 때 ActionSentinel에서 하는 동작
func (sActionSentinel *ActionSentinel) PressNoBtn(s *discordgo.Session, r *discordgo.MessageReaction) {
	// 게임 진행과 관련된 메세지에 달린 리액션 지운다
	if sActionSentinel.filterReaction(s, r) {
		return
	}
}

// PressDirBtn 좌 -1, 우 1 사용자가 좌우 방향 이모티콘을 눌렀을 때 ActionSentinel에서 하는 동작
func (sActionSentinel *ActionSentinel) PressDirBtn(s *discordgo.Session, r *discordgo.MessageReaction, dir int) {
	// 게임 진행과 관련된 메세지에 달린 리액션 지운다
	if sActionSentinel.filterReaction(s, r) {
		return
	}
}

// PressBmkBtn DB에 저장된 정보를 load 하는 동작
func (sActionSentinel *ActionSentinel) PressBmkBtn(s *discordgo.Session, r *discordgo.MessageReaction) {
	//do nothing
}

// InitState 함수는 ActionSentinel state가 시작할 때
// 센티넬 직업을 가진 유저에게 센티넬 동작 DM을 보내고 이를 ActionSentinel 멤버 변수로 저장합니다.
func (sActionSentinel *ActionSentinel) InitState() {
	// sentinel role을 가지고 있는 유저들에게 능력사용 메세지 보낸 후 MessageID 저장
	g := sActionSentinel.g
	g.Session.ChannelMessageEdit(g.ChanID, g.GameStateMID, "플레이어들의 선택 기다리는 중...")
	role := &Sentinel{}
	rIdx := FindRoleIdx(role, sActionSentinel.g.RoleSeq)
	if rIdx != -1 {
		users := g.GetOriRoleUsers(role)
		sActionSentinel.UserChoice = make(chan Choice)
		sActionSentinel.sentinelMsgsID = make(map[string]string)
		if len(users) != 0 {
			for _, user := range users {
				sActionSentinel.sentinelMsgsID[user.UserID] = role.SendUserSelectGuide(user, g, 0)
			}
			cnt := 0
			for input := range sActionSentinel.UserChoice {
				if input.num == -1 {
					tar := &TargetObject{2, "", "", -1}
					role.GenLog(tar, input.user, g)
				} else {
					tar := &TargetObject{2, g.UserList[input.num-1].UserID, "", -1}
					role.Action(tar, input.user, g)
					role.GenLog(tar, input.user, g)
				}
				cnt++
				if cnt == len(users) {
					close(sActionSentinel.UserChoice)
				}
			}
		}
	}
	sActionSentinel.stateFinish()
}

// stateFinish 함수는 sentinel role을 가진 user가 능력사용을 끝내고
// ActionSentinel에서 state가 종료되는 시점에서 호출 됩니다.
// 다음 state인 ActionDoppelganger의 InitState() 함수를 호출합니다.
func (sActionSentinel *ActionSentinel) stateFinish() {
	// ActionDoppelganger의 Msg 변수 초기화
	sActionSentinel.g.CurState = NewActionDoppelganger(sActionSentinel.g)
	sActionSentinel.g.CurState.InitState()
}

// filterReaction 함수는 각 스테이트에서 보낸 메세지에 리액션 했는지 거르는 함수이다.
// 각 스테이트에서 보낸 메세지의 아이디와 리액션이 온 아이디가 동일한지 확인 및
// 메세지에 리액션 한 것을 지워주어야 한다.
func (sActionSentinel *ActionSentinel) filterReaction(s *discordgo.Session, r *discordgo.MessageReaction) bool {
	// 현재 스테이트에서 보낸 메세지에 리액션한 게 아니면 거름
	for _, MsgID := range sActionSentinel.sentinelMsgsID {
		if r.MessageID == MsgID {
			// 메세지에 리약션한 거 지워줌
			s.MessageReactionRemove(r.ChannelID, r.MessageID, r.Emoji.Name, r.UserID)
			return false
		}
	}
	return true
}
