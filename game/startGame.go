package game

import (
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
)

type StartGame struct {
	// state 에서 가지고 있는 game
	g *Game
}

// PressNumBtn 사용자가 숫자 이모티콘을 눌렀을 때 StartGame에서 하는 동작
func (sStartGame *StartGame) PressNumBtn(s *discordgo.Session, r *discordgo.MessageReaction, num int) {
}

// PressDisBtn 사용자가 버려진 카드 이모티콘을 눌렀을 때 StartGame에서 하는 동작
func (sStartGame *StartGame) PressDisBtn(s *discordgo.Session, r *discordgo.MessageReaction) {
}

// PressYesBtn 사용자가 yes 이모티콘을 눌렀을 때 StartGame에서 하는 동작
func (sStartGame *StartGame) PressYesBtn(s *discordgo.Session, r *discordgo.MessageReaction) {
}

// PressNoBtn 사용자가 No 이모티콘을 눌렀을 때 StartGame에서 하는 동작
func (sStartGame *StartGame) PressNoBtn(s *discordgo.Session, r *discordgo.MessageReaction) {
}

// PressDirBtn 좌 -1, 우 1 사용자가 좌우 방향 이모티콘을 눌렀을 때 StartGame에서 하는 동작
func (sStartGame *StartGame) PressDirBtn(s *discordgo.Session, r *discordgo.MessageReaction, dir int) {
}

// PressBmkBtn DB에 저장된 정보를 load 하는 동작
func (sStartGame *StartGame) PressBmkBtn(s *discordgo.Session, r *discordgo.MessageReaction) {
	//do nothing
}

// InitState 함수는 StartGame state가 시작할 때 진짜로 게임이 시작되므로
// game에 UserList에 직업을 랜덤 할당해주고 각 유저에게 직업소개 개인 DM을 보낸 후 센티넬 state를 시작합니다.
func (sStartGame *StartGame) InitState() {
	g := sStartGame.g
	lenuser := len(g.UserList)
	startEmbed := embed.NewEmbed()
	startEmbed.SetTitle("게임 시작!")
	userListMsg := ""
	for _, user := range g.UserList {
		userListMsg += "`" + user.nick + "`\n"
	}
	startEmbed.AddField("유저 목록", userListMsg)
	roleListMsg := ""
	for _, role := range g.RoleView {
		roleListMsg += "`" + role.String() + "`"
		roleListMsg += sStartGame.g.getTeamMark(role.String()) + "\n"
	}
	startEmbed.AddField("설정된 직업 목록", roleListMsg)
	startEmbed.InlineAllFields()
	go g.Session.ChannelMessageSendEmbed(g.ChanID, startEmbed.MessageEmbed)
	// game의 RoleView는 unsorted이기 때문에 섞여도 된다.
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(g.RoleView), func(i, j int) {
		g.RoleView[i], g.RoleView[j] = g.RoleView[j], g.RoleView[i]
	})
	// indexTable은 [len(UserList)][len(RoleSeq)] 만큼의 크기를 갖는다.
	// role이 중복되지 않음에 주의, len(RoleView) > len(RoleSeq)
	g.roleIdxTable = make([][]int, lenuser)
	g.OriRoleIdxTable = make([][]int, lenuser)
	lenrole := len(g.RoleSeq)
	// user에게 랜덤배정되는 role의 수는 lenuser만큼인데
	// indexTable에서는 RoleSeq의 인덱스 번호로 할당해야하고
	// 정작 중복안된 직업이 들어있는건 RoleView이므로
	// RoleView에서 랜덤으로 직업을 뽑아(위에서 suffle로 RoleView를 섞었음)
	// RoleSeq에서 인덱스 위치를 찾아서 indexTable 업데이트 해야함
	for i := 0; i < lenuser; i++ {
		g.roleIdxTable[i] = make([]int, lenrole)
		g.roleIdxTable[i][FindRoleIdx(g.RoleView[i], g.RoleSeq)] = 1
		g.OriRoleIdxTable[i] = make([]int, lenrole)
		g.OriRoleIdxTable[i][FindRoleIdx(g.RoleView[i], g.RoleSeq)] = 1
	}
	// len(UserList) + 3 == len(RoleView) 이므로
	// lenuser 만큼 할당하고 남은 3개는 DisRole에 들어가야 함
	for i := 0; i < 3; i++ {
		g.DisRole[i] = g.RoleView[lenuser+i]
	}
	ch := make(chan bool, len(g.UserList))
	g.Session.ChannelMessageEdit(g.ChanID, g.GameStateMID, "각 직업별 선택지 전송중입니다...")
	// 각 유저에게 직업소개 DM을 보낸다.
	for _, item := range g.UserList {
		go func(item *User, g *Game) {
			userrole := g.GetRole(item.UserID)
			guide := g.RG[userrole.ID()]
			msg := ""
			for _, item := range guide.RoleGuide {
				msg += item + "\n"
			}
			g.Session.ChannelMessageSendEmbed(item.dmChanID, embed.NewGenericEmbed("당신의 직업은 `"+userrole.String()+"` 입니다", msg))
			ch <- true
		}(item, g)
	}
	for i := 0; i < len(g.UserList); i++ {
		<-ch
	}
	// StartGame state 종료 다음 state 시작
	sStartGame.stateFinish()
}

// stateFinish 함수는 StartGame에서 state가 종료되는 시점에서 호출 됩니다.
// 다음 state인 ActionSentinel의 InitState() 함수를 호출합니다.
func (sStartGame *StartGame) stateFinish() {
	sStartGame.g.CurState = &ActionSentinel{sStartGame.g, nil, nil}
	sStartGame.g.CurState.InitState()
}

// filterReaction 함수는 각 스테이트에서 보낸 메세지에 리액션 했는지 거르는 함수이다.
// 각 스테이트에서 보낸 메세지의 아이디와 리액션이 온 아이디가 동일한지 확인 및
// 메세지에 리액션 한 것을 지워주어야 한다.
func (sStartGame *StartGame) filterReaction(s *discordgo.Session, r *discordgo.MessageReaction) bool {
	return false
}
