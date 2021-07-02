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
// game에 UserList에 직업을 랜덤 할당해주고 각 유저에게 직업소개 개인 DM을 보낸 후
// 센티넬 직업을 가진 유저에게 센티넬 동작 DM을 보내고 이를 StartGame 멤버 변수로 저장합니다.
func (sStartGame *StartGame) InitState() {
	g := sStartGame.g
	lenrole := len(g.RoleView)
	lenuser := len(g.UserList)
	randTable := make([]int, lenrole)
	for i := 0; i < lenrole; i++ {
		randTable[i] = i
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(randTable), func(i, j int) {
		randTable[i], randTable[j] = randTable[j], randTable[i]
	})
	idxTable := make([][]int, lenuser)
	for i := 0; i < lenuser; i++ {
		idxTable[i] = make([]int, lenrole)
	}

	for i := 0; i < lenuser; i++ {
		idxTable[i][randTable[i]] = 1
	}
	g.oriRoleIdxTable = idxTable
	g.DisRole[0] = g.RoleView[randTable[lenuser]]
	g.DisRole[1] = g.RoleView[randTable[lenuser+1]]
	g.DisRole[2] = g.RoleView[randTable[lenuser+2]]
	for _, item := range g.UserList {
		go func(item *User, g *Game) {
			userrole := g.GetRole(item.UserID)
			guide := g.RG[userrole.ID()]
			msg := ""
			for _, item := range guide.RoleGuide {
				msg += item + "\n"
			}
			g.Session.ChannelMessageSendEmbed(item.dmChanID, embed.NewGenericEmbed("당신의 직업은 `"+userrole.String()+"` 입니다", msg))
		}(item, g)
	}
}

// stateFinish 함수는 sentinel role을 가진 user가 능력사용을 끝내고
// StartGame에서 state가 종료되는 시점에서 호출 됩니다.
// 다음 state인 ActionDoppelganger의 InitState() 함수를 호출합니다.
func (sStartGame *StartGame) stateFinish(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// do nothing
}

// filterReaction 함수는 각 스테이트에서 보낸 메세지에 리액션 했는지 거르는 함수이다.
// 각 스테이트에서 보낸 메세지의 아이디와 리액션이 온 아이디가 동일한지 확인 및
// 메세지에 리액션 한 것을 지워주어야 한다.
func (sStartGame *StartGame) filterReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// do nothing
}
