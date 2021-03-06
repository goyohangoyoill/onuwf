package game

import (
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
)

// Prepare is test
type StateVote struct {
	g         *Game
	votedList []int
	mostVotes int
	mutex     chan bool
}

func NewStateVote(g *Game) *StateVote {
	ac := &StateVote{}
	ac.g = g
	ac.votedList = make([]int, len(g.UserList))
	ac.mutex = make(chan bool, 1)
	return ac
}

// PressNumBtn 사용자가 숫자 이모티콘을 눌렀을 때 state에서 하는 동작
func (v *StateVote) PressNumBtn(s *discordgo.Session, r *discordgo.MessageReaction, num int) {
	//num를 받음
	//해당 index list count +1
	rUserNum := 9999
	for i := 0; i < num; i++ {
		if r.UserID == v.g.UserList[i].UserID {
			rUserNum = i
			break
		}
	}
	if rUserNum < num {
		num = num + 1
	}
	v.votedList[num-1]++
	s.ChannelMessageDelete(r.ChannelID, r.MessageID)
	v.mutex <- true
	// User.voteUserId에 누구에게 투표했는지 유저ID 저장
	v.setVoteUserId(r.UserID, v.g.UserList[num-1].UserID)
	// 각 유저별 투표내용 로그 작성
	v.setUserVoteLog(v.g.FindUserByUID(r.UserID).nick, v.g.UserList[num-1].nick)
	// 투표완료시 어떤 유저에게 투표했는지 DM으로 메시지를 보내는 함수
	v.sendVoteCompleteMsgToDm(v.g.FindUserByUID(r.UserID), v.g.UserList[num-1].nick)
	if !v.votedAllUsers() {
		<-v.mutex
		return
	}
	<-v.mutex
	v.mostVotes = 0
	for i := 0; i < len(v.g.UserList); i++ {
		if v.mostVotes < v.votedList[i] {
			v.mostVotes = v.votedList[i]
		}
	}
	switch v.mostVotes {
	case 1:
		rMsg := "모두 1표씩 투표받아 아무도 사망하지 않았습니다"
		s.ChannelMessageSendEmbed(v.g.ChanID, embed.NewGenericEmbed("전원 생존", rMsg))
	default:
		rMsg := ""
		for i, user := range v.g.UserList {
			if v.mostVotes == v.votedList[i] {
				if len(rMsg) > 0 {
					s.ChannelMessageSendEmbed(v.g.ChanID, embed.NewGenericEmbed("", "투표 동점자가 있습니다 ..!"))
				}
				v.g.MostVoted = v.g.UserList[i]
				rMsg = user.nick + "님이 총 " + strconv.Itoa(v.mostVotes) + "표로 처형되었습니다"
				s.ChannelMessageSendEmbed(v.g.ChanID, embed.NewGenericEmbed("투표 결과", rMsg))

				s.ChannelMessageSendEmbed(v.g.ChanID, embed.NewGenericEmbed("", "처형된 사람의 직업은...?"))
				for j := 0; j < 3; j++ {
					s.ChannelMessageSend(v.g.ChanID, "...")
					time.Sleep(time.Second)
				}
				executeTitle := "**처형된 사람의 직업은**"
				executeMsg := "`" + v.g.GetRole(v.g.UserList[i].UserID).String() + "`이었습니다!"
				s.ChannelMessageSendEmbed(v.g.ChanID, embed.NewGenericEmbed(executeTitle, executeMsg))
			}
		}
	}
	// 사냥꾼 능력발동 및 로그 작성하는 함수
	v.hunterSkillMsg(s)
	// 각 팀별 승패여부 확인 후 표시
	v.gameResult(s)
	v.stateFinish()
}

// PressDisBtn 사용자가 버려진 카드 이모티콘을 눌렀을 때 state에서 하는 동작
func (v *StateVote) PressDisBtn(_ *discordgo.Session, _ *discordgo.MessageReaction) {
}

// PressYesBtn 사용자가 yes 이모티콘을 눌렀을 때 state에서 하는 동작
func (v *StateVote) PressYesBtn(_ *discordgo.Session, _ *discordgo.MessageReaction) {
}

// PressNoBtn 사용자가 No 이모티콘을 눌렀을 때 state에서 하는 동작
func (v *StateVote) PressNoBtn(_ *discordgo.Session, _ *discordgo.MessageReaction) {
}

// PressDirBtn 좌 -1, 우 1 사용자가 좌우 방향 이모티콘을 눌렀을 때 state에서 하는 동작
func (v *StateVote) PressDirBtn(_ *discordgo.Session, _ *discordgo.MessageReaction, _ int) {
}

// PressBmkBtn DB에 저장된 정보를 load 하는 동작
func (v *StateVote) PressBmkBtn(_ *discordgo.Session, _ *discordgo.MessageReaction) {
	//do nothing
}

// InitState 함수는 스테이트가 시작할 때 필요한 메세지를 생성하고 채널이나 개인DM으로 메세지를 보낸 후
// 메세지 객체를 스테이트의 멤버로 저장합니다.
// 이 함수는 이전 스테이트가 끝나는 시점에 호출되어야 합니다.
func (v *StateVote) InitState() {
	delaySec := v.g.config.VoteDelaySec
	if delaySec > 0 {
		msg := strconv.Itoa(delaySec) + "초 후에 투표가 시작됩니다"
		v.g.Session.ChannelMessageSend(v.g.ChanID, msg)
		time.Sleep(time.Duration(delaySec) * time.Second)
	}
	msg := ""
	v.g.AppendLog(msg)
	v.g.Session.ChannelMessageEdit(v.g.ChanID, v.g.GameStateMID, "투표용지 전달중...")
	VoteProcess(v.g.Session, v.g)
	v.g.Session.ChannelMessageEdit(v.g.ChanID, v.g.GameStateMID, "투표 진행중...")
	v.g.AppendLog("투표 결과")
}

// stateFinish 함수는 현재 state가 끝나고 다음 state로 넘어갈 때 호출되는 함수입니다.
// game의 CurState 변수에 다음 state를 생성해서 할당해준 다음
// 다음 state의 InitState() 함수를 이 함수 안에서 호출해야 합니다
func (v *StateVote) stateFinish() {
	v.g.SendLogMsg(v.g.ChanID)
	v.g.Session.ChannelMessageEdit(v.g.ChanID, v.g.GameStateMID, "게임 종료.")
	v.g.GameStartedChan <- false
}

// filterReaction 함수는 각 스테이트에서 보낸 메세지에 리액션 했는지 거르는 함수이다.
// 각 스테이트에서 보낸 메세지의 아이디와 리액션이 온 아이디가 동일한지 확인 및
// 메세지에 리액션 한 것을 지워주어야 한다.
func (v *StateVote) filterReaction(_ *discordgo.Session, _ *discordgo.MessageReaction) bool {
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
	voteEmbed.SetDescription("처형하실 대상에게 투표해주세요")
	num := len(g.UserList)
	for i := 0; i < num-1; i++ {
		j := i
		if j >= UserNum {
			j = j + 1
		}
		voteEmbed.AddField(strconv.Itoa(i+1)+"번 ", g.UserList[j].nick)
	}
	voteEmbed.SetAuthor(g.UserList[UserNum].title + " " + "`" + g.UserList[UserNum].nick + "`")
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
}

// User.voteUserId에 누구에게 투표했는지 유저ID 저장
func (v *StateVote) setVoteUserId(voteUserId, votedUserId string) {
	for _, user := range v.g.UserList {
		if user.UserID == voteUserId {
			user.voteUserId = votedUserId
			break
		}
	}
}

// 해당 직업이 처형되었는지 확인하는 함수
func (v *StateVote) isExecutedRole(i int, role string) bool {
	if v.mostVotes == v.votedList[i] && v.mostVotes > 1 {
		if v.g.GetRole(v.g.UserList[i].UserID).String() == role {
			return true
		}
	}
	return false
}

// 사냥꾼 능력발동 및 로그 작성하는 함수
func (v *StateVote) hunterSkillMsg(s *discordgo.Session) {
	for i, user := range v.g.UserList {
		if v.isExecutedRole(i, "사냥꾼") {
			hunterMsg := "`" + user.Nick() + "`가 `" + v.g.FindUserByUID(user.voteUserId).nick + "` 사냥에 성공했습니다..!\n이 사람의 직업은..?"
			s.ChannelMessageSendEmbed(v.g.ChanID, embed.NewGenericEmbed("사냥꾼 능력발동!", hunterMsg))
			for j := 0; j < 3; j++ {
				s.ChannelMessageSend(v.g.ChanID, "...")
				time.Sleep(time.Second)
			}
			hunterMsg = "`" + v.g.GetRole(user.voteUserId).String() + "`이었습니다!"
			s.ChannelMessageSendEmbed(v.g.ChanID, embed.NewGenericEmbed("", hunterMsg))
			break
		}
	}
}

// 투표완료시 어떤 유저에게 투표했는지 DM으로 메시지를 보내는 함수
func (v *StateVote) sendVoteCompleteMsgToDm(voteUser *User, votedUserNick string) {
	title := "투표 완료"
	msg := "`" + votedUserNick + "`에게 투표하셨습니다"
	v.g.Session.ChannelMessageSendEmbed(voteUser.dmChanID, embed.NewGenericEmbed(title, msg))
	v.g.Session.ChannelMessageSend(v.g.ChanID, "`"+voteUser.nick+"` 님이 투표하셨습니다.")
}

// 각 유저별 투표 내용 작성
func (v *StateVote) setUserVoteLog(from, to string) {
	msg := "`" + from + "` -> `" + to + "`"
	v.g.AppendLog(msg)
}

// 모든 유저가 투표했는지 확인
func (v *StateVote) votedAllUsers() bool {
	for _, user := range v.g.UserList {
		if len(user.voteUserId) == 0 {
			return false
		}
	}
	return true
}

// 처형당한 직업이 있는지 확인
func (v *StateVote) existsExecutedRole(roleName string) bool {
	for i := 0; i < len(v.g.UserList); i++ {
		if v.isExecutedRole(i, roleName) {
			return true
		}
	}
	return false
}

// 각 팀별 승패여부 확인 후 표시
func (v *StateVote) gameResult(s *discordgo.Session) {
	voteResultEmbed := embed.NewEmbed()
	voteResultEmbed.SetTitle("게임 결과")
	// 마을주민팀 승패여부
	if v.g.userExistsOfThisTeam("Villager") {
		v.g.VillagerTeamWin = false
		switch len(v.g.GetRoleUsers(&Werewolf{})) {
		case 0:
			switch v.mostVotes {
			case 1:
				v.g.VillagerTeamWin = true
				voteResultEmbed.AddField("마을주민팀 승리", "모두 1표씩 투표받아 아무도 사망하지 않았습니다.")
			default:
				voteResultEmbed.AddField("마을주민팀 패배", "누군가 사망했습니다.")
			}
		default:
			switch v.existsExecutedRole("늑대인간") {
			case true:
				v.g.VillagerTeamWin = true
				voteResultEmbed.AddField("마을주민팀 승리", "늑대인간이 처형되었습니다.")
			case false:
				voteResultEmbed.AddField("마을주민팀 패배", "늑대인간이 처형되지 않았습니다.")
			}
		}
	}
	// 늑대인간팀 승패여부
	if v.g.userExistsOfThisTeam("Werewolf") {
		v.g.WerewolfTeamWin = true
		switch len(v.g.GetRoleUsers(&Werewolf{})) {
		case 0:
			switch {
			case v.mostVotes == 1:
				v.g.WerewolfTeamWin = false
				voteResultEmbed.AddField("늑대인간팀 패배", "아무도 사망하지 않았습니다.")
			case v.existsExecutedRole("하수인"):
				v.g.WerewolfTeamWin = false
				voteResultEmbed.AddField("늑대인간팀 패배", "하수인이 처형되었습니다.")
			case v.g.userExistsOfThisTeam("Tanner") && v.existsExecutedRole("무두장이"):
				v.g.WerewolfTeamWin = false
				voteResultEmbed.AddField("늑대인간팀 패배", "무두장이가 처형되었습니다.")
			default:
				voteResultEmbed.AddField("늑대인간팀 승리", "하수인, 무두장이가 살아남았습니다.")
			}
		default:
			switch {
			case v.existsExecutedRole("늑대인간"):
				v.g.WerewolfTeamWin = false
				voteResultEmbed.AddField("늑대인간팀 패배", "늑대인간이 처형되었습니다.")
			case v.g.userExistsOfThisTeam("Tanner") && v.existsExecutedRole("무두장이"):
				v.g.WerewolfTeamWin = false
				voteResultEmbed.AddField("늑대인간팀 패배", "무두장이가 처형되었습니다.")
			default:
				voteResultEmbed.AddField("늑대인간팀 승리", "늑대인간이 살아남았습니다.")
			}
		}
	}
	// 무두장이팀 승패여부
	if v.g.userExistsOfThisTeam("Tanner") {
		v.g.TannerTeamWin = false
		switch v.existsExecutedRole("무두장이") {
		case true:
			v.g.TannerTeamWin = true
			voteResultEmbed.AddField("무두장이팀 승리", "무두장이가 처형되었습니다.")
		case false:
			voteResultEmbed.AddField("무두장이팀 패배", "무두장이가 살아남았습니다.")
		}
	}
	s.ChannelMessageSendEmbed(v.g.ChanID, voteResultEmbed.MessageEmbed)
}
