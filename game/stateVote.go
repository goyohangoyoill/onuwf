package game

import (
	"strconv"

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

// PressBmkBtn DB에 저장된 정보를 load 하는 동작
func (v *StateVote) PressBmkBtn(s *discordgo.Session, r *discordgo.MessageReaction) {
	//do nothing
}

// InitState 함수는 스테이트가 시작할 때 필요한 메세지를 생성하고 채널이나 개인DM으로 메세지를 보낸 후
// 메세지 객체를 스테이트의 멤버로 저장합니다.
// 이 함수는 이전 스테이트가 끝나는 시점에 호출되어야 합니다.
func (v *StateVote) InitState() {
	//	v.G.UserList = append(v.G.UserList, NewUser(v.G.MasterID, "juhur", v.G.ChanID, v.G.ChanID))
	//v.G.UserList = append(v.G.UserList, NewUser(v.G.MasterID, "kalee", v.G.ChanID, v.G.ChanID))
	//v.G.UserList = append(v.G.UserList, NewUser(v.G.MasterID, "min-jo", v.G.ChanID, v.G.ChanID))
	msg := ""
	v.G.AppendLog(msg)
	VoteProcess(v.G.Session, v.G)
}

// stateFinish 함수는 현재 state가 끝나고 다음 state로 넘어갈 때 호출되는 함수입니다.
// game의 CurState 변수에 다음 state를 생성해서 할당해준 다음
// 다음 state의 InitState() 함수를 이 함수 안에서 호출해야 합니다
func (v *StateVote) stateFinish() {

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
	voteEmbed.SetAuthor(g.UserList[UserNum].title + " " + g.UserList[UserNum].nick)
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
