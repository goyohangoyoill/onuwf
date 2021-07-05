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
func (sStartGame *StartGame) PressNumBtn(s *discordgo.Session, r *discordgo.MessageReactionAdd, num int) {
	// do nothing
}

// PressDisBtn 사용자가 버려진 카드 이모티콘을 눌렀을 때 StartGame에서 하는 동작
func (sStartGame *StartGame) PressDisBtn(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// do nothing
}

// PressYesBtn 사용자가 yes 이모티콘을 눌렀을 때 StartGame에서 하는 동작
func (sStartGame *StartGame) PressYesBtn(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// do nothing
}

// PressNoBtn 사용자가 No 이모티콘을 눌렀을 때 StartGame에서 하는 동작
func (sStartGame *StartGame) PressNoBtn(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// do nothing
}

// PressDirBtn 좌 -1, 우 1 사용자가 좌우 방향 이모티콘을 눌렀을 때 StartGame에서 하는 동작
func (sStartGame *StartGame) PressDirBtn(s *discordgo.Session, r *discordgo.MessageReactionAdd, dir int) {
	// do nothing
}

// InitState 함수는 StartGame state가 시작할 때 진짜로 게임이 시작되므로
// game에 UserList에 직업을 랜덤 할당해주고 각 유저에게 직업소개 개인 DM을 보냅니다.
func (sStartGame *StartGame) InitState() {
	g := sStartGame.g
	lenuser := len(g.UserList)
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(g.RoleView), func(i, j int) {
		g.RoleView[i], g.RoleView[j] = g.RoleView[j], g.RoleView[i]
	})
	g.roleIdxTable = make([][]int, lenuser)
	g.oriRoleIdxTable = make([][]int, lenuser)
	lenrole := len(g.RoleSeq)
	for i := 0; i < lenuser; i++ {
		g.roleIdxTable[i] = make([]int, lenrole)
		g.roleIdxTable[i][FindRoleIdx(g.RoleView[i], g.RoleSeq)] = 1
		g.oriRoleIdxTable[i] = make([]int, lenrole)
		g.oriRoleIdxTable[i][FindRoleIdx(g.RoleView[i], g.RoleSeq)] = 1
	}
	for i := 0; i < 3; i++ {
		g.DisRole[i] = g.RoleView[lenuser+i]
	}
	ch := make(chan bool, len(g.UserList))
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
	sStartGame.stateFinish(g.Session, nil)
}

// stateFinish 함수는 개인별 안내 메시지가 모두 전송되고,
// StartGame에서 state가 종료되는 시점에서 호출 됩니다.
// 다음 state인 ActionDoppelganger의 InitState() 함수를 호출합니다.
func (sStartGame *StartGame) stateFinish(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	sStartGame.g.CurState = &ActionSentinel{sStartGame.g, nil}
	s.ChannelMessageSend(sStartGame.g.ChanID, "모두에게 직업이 배정되었습니다.")
	sStartGame.g.CurState.InitState()
}

// filterReaction 함수는 각 스테이트에서 보낸 메세지에 리액션 했는지 거르는 함수이다.
// 각 스테이트에서 보낸 메세지의 아이디와 리액션이 온 아이디가 동일한지 확인 및
// 메세지에 리액션 한 것을 지워주어야 한다.
func (sStartGame *StartGame) filterReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// do nothing
}
