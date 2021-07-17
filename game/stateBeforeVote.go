package game

import (
	"github.com/bwmarrin/discordgo"
)

// State 는 리액션 입력이 발생했을 때 현재 상태에 따라 다른 함수를 호출하는 작업을 수행하는
// 인터페이스로, 숫자 버튼, 쓰레기통 버튼, 예/아니오 버튼, 방향 버튼을 인식할 수 있게 구현한다.
type StateBeforeVote struct {
	G    *Game
	Info map[string]*DMInfo
}

func NewStateBeforeVote(g *Game) *StateBeforeVote {
	ac := &StateBeforeVote{}
	ac.G = g
	ac.Info = make(map[string]*DMInfo)
	return ac
}

// 사용자 인원수 3 ~ 26
// num: 0 ~ 23
// PressNumBtn 사용자가 숫자 이모티콘을 눌렀을 때 state에서 하는 동작
func (sb *StateBeforeVote) PressNumBtn(s *discordgo.Session, r *discordgo.MessageReaction, num int) {
	player := sb.G.FindUserByUID(r.UserID)
	curInfo := sb.Info[player.UserID]
	s.ChannelMessageDelete(r.ChannelID, curInfo.MsgID)

	curInfo.Choice <- num
}

// PressDisBtn 사용자가 버려진 카드 이모티콘을 눌렀을 때 state에서 하는 동작
func (sb *StateBeforeVote) PressDisBtn(s *discordgo.Session, r *discordgo.MessageReaction) {
}

// PressYesBtn 사용자가 yes 이모티콘을 눌렀을 때 state에서 하는 동작
func (sb *StateBeforeVote) PressYesBtn(s *discordgo.Session, r *discordgo.MessageReaction) {
}

// PressNoBtn 사용자가 No 이모티콘을 눌렀을 때 state에서 하는 동작
func (sb *StateBeforeVote) PressNoBtn(s *discordgo.Session, r *discordgo.MessageReaction) {
}

// PressDirBtn 좌 -1, 우 1 사용자가 좌우 방향 이모티콘을 눌렀을 때 state에서 하는 동작
func (sb *StateBeforeVote) PressDirBtn(s *discordgo.Session, r *discordgo.MessageReaction, dir int) {
}

// InitState 함수는 스테이트가 시작할 때 필요한 메세지를 생성하고 채널이나 개인DM으로 메세지를 보낸 후
// 메세지 객체를 스테이트의 멤버로 저장합니다.
// 이 함수는 이전 스테이트가 끝나는 시점에 호출되어야 합니다.
func (sb *StateBeforeVote) InitState() {
	//불면증, 주정뱅이 각각의 rolemsg는 이전 state에서 보내야함
	//능력사용

	// 불면증환자
	role := GenerateRole(9)
	rIdx := FindRoleIdx(role, sb.G.RoleSeq)
	if rIdx != -1 {
		InsomUsers := sb.G.GetOriRoleUsers(role)
		for i := 0; i < len(InsomUsers); i++ {
			//role := GenerateRole(15)
			tar := &TargetObject{2, InsomUsers[i].UserID, "", 0}
			role.SendUserSelectGuide(InsomUsers[i], sb.G, 0)
			role.Action(tar, InsomUsers[i], sb.G)
			role.GenLog(tar, InsomUsers[i], sb.G)
		}
	}
	// 주정뱅이
	role = GenerateRole(8)
	rIdx = FindRoleIdx(role, sb.G.RoleSeq)
	if rIdx != -1 {
		DrunkUsers := sb.G.GetOriRoleUsersWithoutDpl(role)
		for _, user := range DrunkUsers {
			curInfo := &DMInfo{"", make(chan int), 0}
			sb.Info[user.UserID] = curInfo
			curInfo.MsgID = role.SendUserSelectGuide(user, sb.G, 1)
		}
		//curInfo := &DMInfo{"", make(chan int), 0}
		curInfo := sb.Info
		for i := 0; i < len(DrunkUsers); i++ {
			//curInfo.MsgID = role.SendUserSelectGuide(DrunkUsers[i], sb.G, 1)
			//curInfo := sb.Info
			input := <-curInfo[DrunkUsers[i].UserID].Choice
			tar := &TargetObject{1, DrunkUsers[i].UserID, "", input - 1}

			role.Action(tar, DrunkUsers[i], sb.G)
			role.GenLog(tar, DrunkUsers[i], sb.G)
		}
		sb.G.AppendLog("")
	}
	sb.stateFinish()
}

// stateFinish 함수는 현재 state가 끝나고 다음 state로 넘어갈 때 호출되는 함수입니다.
// game의 CurState 변수에 다음 state를 생성해서 할당해준 다음
// 다음 state의 InitState() 함수를 이 함수 안에서 호출해야 합니다
func (sb *StateBeforeVote) stateFinish() {
	sb.G.CurState = NewStateVote(sb.G)
	sb.G.CurState.InitState()
}

// filterReaction 함수는 각 스테이트에서 보낸 메세지에 리액션 했는지 거르는 함수이다.
// 각 스테이트에서 보낸 메세지의 아이디와 리액션이 온 아이디가 동일한지 확인 및
// 메세지에 리액션 한 것을 지워주어야 한다.
func (sb *StateBeforeVote) filterReaction(s *discordgo.Session, r *discordgo.MessageReaction) bool {
	return false
}
