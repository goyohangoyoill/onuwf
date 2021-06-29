/* onuwf 는 보드게임 "한밤의 늑대인간" 을 디스코드 봇으로 구현하는 프로젝트입니다. */

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	wfGame "onuwf.com/game"
	wfUtil "onuwf.com/util"

	"github.com/bwmarrin/discordgo"
)

var (
	isUserIn      map[string]bool
	isGuildChanIn map[string]bool
	uidToGameData map[string]*wfGame.Game
)

func init() {
	wfUtil.EnvInit()
	wfUtil.RoleGuideInit()
	wfUtil.LoggerInit()

	isUserIn = make(map[string]bool)
	isGuildChanIn = make(map[string]bool)
	uidToGameData = make(map[string]*wfGame.Game)
}

func main() {
	dg, err := discordgo.New("Bot " + env["dgToken"])
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	dg.AddHandler(messageCreate)
	dg.AddHandler(messageReactionAdd)
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	dg.Close()
}

func startgame(s *discordgo.Session, m *discordgo.MessageCreate) {

}

// messageCreate() 입력한 메시지를 처리하는 함수
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "ㅁ시작" {
		if isUserIn[m.Author.ID] {
			s.ChannelMessageSend(m.ChannelID, "게임을 진행중인 사용자입니다.")
			return
		}
		if isGuildChanIn[m.ChannelID+m.GuildID] {
			s.ChannelMessageSend(m.ChannelID, "게임을 진행중인 채널입니다.")
			return
		}
		isGuildChanIn[m.ChannelID+m.GuildID] = true
		isUserIn[m.Author.ID] = true
		go startgame(s, m)
	}
	if m.Content == "ㅁ강제종료" {
		if isInGame[m.GuildID+m.ChannelID] {
			isInGame[m.GuildID+m.ChannelID] = false
		}
	}
}

// messageReactionAdd 함수는 인게임 버튼 이모지 상호작용 처리를 위한 이벤트 핸들러 함수입니다.
func messageReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// 봇 자기자신의 리액션 무시.
	if r.UserID == s.State.User.ID {
		return
	}

	// 게임 참가중이 아닌 사용자의 리액션 무시.
	if isUserIn[r.UserID] {
		return
	}

	g := uidToGameData[r.UserID]
	// 숫자 이모지 선택.
	for i := 1; i < 10; i++ {
		var ch rune
		ch = '0' + rune(i)
		emjID := "n" + string(ch)
		if r.Emoji.Name == emj[emjID] {
			g.curState.PressNumBtn(s, r, i)
		}
	}
	switch r.Emoji.Name {
	case emj["DISCARD"]:
		// 쓰레기통 이모지 선택.
		go g.curState.PressDisBtn(s, r)
	case emj["YES"]:
		// O 이모지 선택.
		go g.curState.PressYesBtn(s, r)
	case emj["NO"]:
		// X 이모지 선택.
		go g.curState.PressNoBtn(s, r)
	case emj["LEFT"]:
		// 왼쪽 화살표 선택.
		go g.curState.PressDirBtn(s, r, -1)
	case emj["RIGHT"]:
		// 오른쪽 화살표 선택.
		go g.curState.PressDirBtn(s, r, 1)
	}
	if r.GuildID == curGame.guildID && r.ChannelID == curGame.chanID && (r.MessageID == curGame.enterGameMsgID || r.MessageID == curGame.roleAddMsgID) {
		s.MessageReactionRemove(r.ChannelID, r.MessageID, r.Emoji.Name, r.UserID)
	}

}
