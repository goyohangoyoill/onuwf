package game

import (
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
)

// Prepare is test
type StateVote struct {
	G          *Game
	Voted_list []int
	User_num   int
	Vote_count int
}

func NewStateVote(g *Game) *StateVote {
	ac := &StateVote{}
	ac.G = g
	ac.Voted_list = make([]int, len(g.UserList))
	ac.User_num = len(g.UserList)
	ac.Vote_count = 0

	return ac
}

// PressNumBtn 사용자가 숫자 이모티콘을 눌렀을 때 state에서 하는 동작
func (v *StateVote) PressNumBtn(s *discordgo.Session, r *discordgo.MessageReaction, num int) {
	//num를 받음
	//해당 index list count +1

	rUserNum := 9999
	for i := 0; i < num; i++ {
		if r.UserID == v.G.UserList[i].UserID {
			rUserNum = i
			break
		}
	}

	if rUserNum < num {
		num = num + 1
	}
	v.Voted_list[num-1]++
	msg := ""
	msg += v.G.GetRole(r.UserID).String() + " " + v.G.FindUserByUID(r.UserID).nick + " 는 "
	msg += v.G.GetRole(v.G.UserList[num-1].UserID).String() + " " + v.G.UserList[num-1].nick + "에게 투표하였습니다"
	v.G.AppendLog(msg)
	s.ChannelMessageDelete(r.ChannelID, r.MessageID)
	// 투표 내용 저장
	v.setUserVoteId(r.UserID, v.G.UserList[num-1].UserID)
	// 투표 완료 DM메시지
	v.sendVoteCompleteMsgToDm(v.G.FindUserByUID(r.UserID), v.G.UserList[num-1].nick)
	v.Vote_count++
	if v.Vote_count == v.User_num {
		max_value := 0
		for i := 0; i < v.User_num; i++ {
			if max_value < v.Voted_list[i] {
				max_value = v.Voted_list[i]
			}
		}
		voteResultEmbed := embed.NewEmbed()
		voteResultEmbed.SetTitle("투표 결과")
		for i := 0; i < v.User_num; i++ {
			rMsg := ""
			if max_value == v.Voted_list[i] {
				voteResultEmbed.AddField(v.G.UserList[i].nick, "`"+v.G.GetRole(v.G.UserList[i].UserID).String()+"` "+v.G.UserList[i].nick+"는 투표로 사망하였습니다.")
				rMsg += v.G.UserList[i].nick + "는 " + strconv.Itoa(max_value) + "회 지목당해 투표로 사망하였습니다"
				v.G.AppendLog(rMsg)
			}
		}
		s.ChannelMessageSendEmbed(v.G.ChanID, voteResultEmbed.MessageEmbed)
		// 헌터 능력 발동
		v.hunterSkillMsg(s, max_value)
		v.stateFinish()
	}
}

// PressDisBtn 사용자가 버려진 카드 이모티콘을 눌렀을 때 state에서 하는 동작
func (v *StateVote) PressDisBtn(s *discordgo.Session, r *discordgo.MessageReaction) {
	//do nothing
}

// PressYesBtn 사용자가 yes 이모티콘을 눌렀을 때 state에서 하는 동작
func (v *StateVote) PressYesBtn(s *discordgo.Session, r *discordgo.MessageReaction) {
	//do nothing
}

// PressNoBtn 사용자가 No 이모티콘을 눌렀을 때 state에서 하는 동작
func (v *StateVote) PressNoBtn(s *discordgo.Session, r *discordgo.MessageReaction) {
	//do nothing
}

// PressDirBtn 좌 -1, 우 1 사용자가 좌우 방향 이모티콘을 눌렀을 때 state에서 하는 동작
func (v *StateVote) PressDirBtn(s *discordgo.Session, r *discordgo.MessageReaction, dir int) {
	//do nothing
}

// InitState 함수는 스테이트가 시작할 때 필요한 메세지를 생성하고 채널이나 개인DM으로 메세지를 보낸 후
// 메세지 객체를 스테이트의 멤버로 저장합니다.
// 이 함수는 이전 스테이트가 끝나는 시점에 호출되어야 합니다.
func (v *StateVote) InitState() {
	delaySec := v.G.config.VoteDelaySec
	if delaySec > 0 {
		msg := strconv.Itoa(delaySec) + "초 후에 투표가 시작됩니다"
		v.G.Session.ChannelMessageSend(v.G.ChanID, msg)
		time.Sleep(time.Duration(delaySec) * time.Second)
	}
	msg := ""
	v.G.AppendLog(msg)
	v.G.Session.ChannelMessageEdit(v.G.ChanID, v.G.GameStateMID, "투표용지 전달중...")
	VoteProcess(v.G.Session, v.G)
	v.G.Session.ChannelMessageEdit(v.G.ChanID, v.G.GameStateMID, "투표 진행중...")
}

// stateFinish 함수는 현재 state가 끝나고 다음 state로 넘어갈 때 호출되는 함수입니다.
// game의 CurState 변수에 다음 state를 생성해서 할당해준 다음
// 다음 state의 InitState() 함수를 이 함수 안에서 호출해야 합니다
func (v *StateVote) stateFinish() {
	v.G.SendLogMsg(v.G.ChanID)
	v.G.Session.ChannelMessageEdit(v.G.ChanID, v.G.GameStateMID, "게임 종료.")
	v.G.GameStartedChan <- false
}

// filterReaction 함수는 각 스테이트에서 보낸 메세지에 리액션 했는지 거르는 함수이다.
// 각 스테이트에서 보낸 메세지의 아이디와 리액션이 온 아이디가 동일한지 확인 및
// 메세지에 리액션 한 것을 지워주어야 한다.
func (v *StateVote) filterReaction(s *discordgo.Session, r *discordgo.MessageReaction) bool {
	// do nothing
	return false
}

func VoteProcess(s *discordgo.Session, g *Game) {
	//send msg
	//개별로 각 채널에서 수행하게 해야함
	//user vote event handler
	//참가자마다 입력을 처리 *go routine 입력을 state vote preesnumbtn 이용해서 count
	// 결과 값을 visualization (통합채널)
	num := len(g.UserList)
	for i := 0; i < num; i++ {
		go SendVoteDM(s, g, i)
	}

}

func SendVoteDM(s *discordgo.Session, g *Game, UserNum int) {
	voteEmbed := embed.NewEmbed()
	voteEmbed.SetTitle("투표")
	voteEmbed.SetDescription("늑대인간으로 의심되는 대상에게 투표해주세요")
	num := len(g.UserList)
	for i := 0; i < num-1; i++ {
		//이후에 본인 빼도록 수정해야함
		j := i
		if j >= UserNum {
			j = j + 1
		}
		voteEmbed.AddField(strconv.Itoa(i+1)+"번 ", g.UserList[j].nick)
	}
	voteEmbed.SetAuthor(g.UserList[UserNum].nick)
	voteEmbed.InlineAllFields()
	UserDM, _ := s.UserChannelCreate(g.UserList[UserNum].UserID) //g.UserList[0] -> g.UsrList[UserNum] change need(test용)
	voteEmbed.InlineAllFields()
	voteMsg, _ := s.ChannelMessageSendEmbed(UserDM.ID, voteEmbed.MessageEmbed)
	addNumAddEmoji(s, voteMsg, g)
}

func addNumAddEmoji(s *discordgo.Session, msg *discordgo.Message, g *Game) {
	num := len(g.UserList)
	for i := 0; i < num-1; i++ {
		s.MessageReactionAdd(msg.ChannelID, msg.ID, g.Emj["n"+strconv.Itoa(i+1)])
	}
	//s.MessageReactionAdd(msg.ChannelID, msg.ID, g.Emj["n2"])
	//s.MessageReactionAdd(msg.ChannelID, msg.ID, g.Emj["n3"])
	//s.MessageReactionAdd(msg.ChannelID, msg.ID, g.Emj["n4"])
}

// User.voteUserId에 투표내용 저장
func (v *StateVote) setUserVoteId(voteUserId, votedUserId string) {
	for _, user := range v.G.UserList {
		if user.UserID == voteUserId {
			user.voteUserId = votedUserId
			break
		}
	}
}

// 사냥꾼 능력발동 조건을 확인하는 함수
func (v *StateVote) chkUseHunterSkill(i, max_value int) bool {
	if max_value == v.Voted_list[i] {
		if v.G.GetRole(v.G.UserList[i].UserID).String() == "사냥꾼" {
			return true
		}
	}
	return false
}

// 사냥꾼 능력발동 및 로그 작성하는 함수
func (v *StateVote) hunterSkillMsg(s *discordgo.Session, max_value int) {
	for i, user := range v.G.UserList {
		if v.chkUseHunterSkill(i, max_value) {
			hunterTitle := "사냥꾼 능력발동!"
			hunterMsg := "`사냥꾼` `" + user.nick + "`이 "
			hunterMsg += "`" + v.G.GetRole(user.voteUserId).String() + "` "
			hunterMsg += "`" + v.G.FindUserByUID(user.voteUserId).nick + "`를 지목하여 사냥에 성공했습니다!"
			s.ChannelMessageSendEmbed(v.G.ChanID, embed.NewGenericEmbed(hunterTitle, hunterMsg))
			v.G.AppendLog(hunterMsg)
			break
		}
	}
}

// 투표완료시 어떤 유저에게 투표했는지 DM으로 메시지를 보내는 함수
func (v *StateVote) sendVoteCompleteMsgToDm(voteUser *User, votedUserNick string) {
	title := "투표 완료"
	msg := "`" + votedUserNick + "`에게 투표하셨습니다"
	v.G.Session.ChannelMessageSendEmbed(voteUser.dmChanID, embed.NewGenericEmbed(title, msg))
	v.G.Session.ChannelMessageSend(v.G.ChanID, "`"+voteUser.nick+"` 님이 투표하셨습니다.")
}
