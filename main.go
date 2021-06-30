/* onuwf 는 보드게임 "한밤의 늑대인간" 을 디스코드 봇으로 구현하는 프로젝트입니다. */

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	wfGame "onuwf.com/game"

	"github.com/bwmarrin/discordgo"
)

var (
	isUserIn      map[string]bool
	isGuildChanIn map[string]bool
	uidToGameData map[string]*wfGame.Game

	env map[string]string
	emj map[string]string
	rg  []wfGame.RoleGuide
)

func init() {
	env = EnvInit()
	emj = EmojiInit()
	RoleGuideInit(&rg)

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
	if !isGuildChanIn[m.ChannelID] {
		userIDChan := make(chan string)
		curGame := wfGame.NewGame(m.GuildID, m.ChannelID, m.Author.ID, s, rg, emj, userIDChan)
		uidToGameData[m.Author.ID] = curGame
		isGuildChanIn[m.ChannelID] = true
		for {
			// Mutex 필요할 것으로 예상됨.
			curUID := <-userIDChan
			isUserIn[curUID] = true
			uidToGameData[curUID] = curGame
		}
	}
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
	if m.Content == "!test" {
		str := rg[3].RoleGuide[0]
		s.ChannelMessageSend(m.ChannelID, str)
	}
	// if m.Content == "ㅁ투표" {
	// 	thisGame := wfGame.NewGame(m.GuildID, m.ChannelID, m.Author.ID, rg, nil)
	// 	thisGame.UserList = append(thisGame.UserList, &wfGame.User{m.Author.ID, "aaa", m.ChannelID, m.ChannelID}
	// 	thisGame.UserList = append(thisGame.UserList, &wfGame.User{m.Author.ID, "bbb", m.ChannelID, m.ChannelID}
	// 	thisGame.UserList = append(thisGame.UserList, &wfGame.User{m.Author.ID, "ccc", m.ChannelID, m.ChannelID}
	// 	thisGame.UserList = append(thisGame.UserList, &wfGame.User{m.Author.ID, "ddd", m.ChannelID, m.ChannelID}
	// 	thisGame.UserList = append(thisGame.UserList, &wfGame.User{m.Author.ID, "eee", m.ChannelID, m.ChannelID}
	// 	thisGame.UserList = append(thisGame.UserList, &wfGame.User{m.Author.ID, "fff", m.ChannelID, m.ChannelID}
	// 	thisGame.CurState = wfGame.StateVote

	// 	wfGame.VoteProcess()
	// }
}

// messageReactionAdd 함수는 인게임 버튼 이모지 상호작용 처리를 위한 이벤트 핸들러 함수입니다.
func messageReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// 봇 자기자신의 리액션 무시.
	if r.UserID == s.State.User.ID {
		return
	}

	// 게임 참가중이 아닌 사용자의 리액션 무시.
	if !isUserIn[r.UserID] {
		return
	}

	g := uidToGameData[r.UserID]
	// 숫자 이모지 선택.
	for i := 1; i < 10; i++ {
		var ch rune
		ch = '0' + rune(i)
		emjID := "n" + string(ch)
		if r.Emoji.Name == emj[emjID] {
			g.CurState.PressNumBtn(s, r, i)
		}
	}
	switch r.Emoji.Name {
	case emj["DISCARD"]:
		// 쓰레기통 이모지 선택.
		go g.CurState.PressDisBtn(s, r)
	case emj["YES"]:
		// O 이모지 선택.
		go g.CurState.PressYesBtn(s, r)
	case emj["NO"]:
		// X 이모지 선택.
		go g.CurState.PressNoBtn(s, r)
	case emj["LEFT"]:
		// 왼쪽 화살표 선택.
		go g.CurState.PressDirBtn(s, r, -1)
	case emj["RIGHT"]:
		// 오른쪽 화살표 선택.
		go g.CurState.PressDirBtn(s, r, 1)
	}
	if r.GuildID == g.GuildID && r.ChannelID == g.ChanID {
		s.MessageReactionRemove(r.ChannelID, r.MessageID, r.Emoji.Name, r.UserID)
	}

}

// EnvInit 설치 환경 불러오기.
func EnvInit() map[string]string {
	envFile, err := os.Open("asset/env.json")
	if err != nil {
		log.Fatal(err)
	}
	defer envFile.Close()

	var byteValue []byte
	byteValue, err = ioutil.ReadAll(envFile)
	if err != nil {
		log.Fatal(err)
	}
	env := make(map[string]string)
	json.Unmarshal([]byte(byteValue), &env)
	return env
}

// RoleGuideInit 직업 가이드 에셋 불러오기.
func RoleGuideInit(rg *[]wfGame.RoleGuide) {
	rgFile, err := os.Open("asset/role_guide.json")
	if err != nil {
		log.Fatal(err)
	}
	defer rgFile.Close()

	var byteValue []byte
	byteValue, err = ioutil.ReadAll(rgFile)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal([]byte(byteValue), rg)
}

// EmojiInit 이모지 맵에 불러오기.
func EmojiInit() map[string]string {
	emjFile, err := os.Open("asset/emoji.json")
	if err != nil {
		log.Fatal(err)
	}
	defer emjFile.Close()

	var byteValue []byte
	byteValue, err = ioutil.ReadAll(emjFile)
	if err != nil {
		log.Fatal(err)
	}
	emj := make(map[string]string)
	json.Unmarshal([]byte(byteValue), &emj)
	return emj
}
