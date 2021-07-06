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
	"time"

	wfGame "onuwf.com/game"

	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
)

const (
	prefix = "ㅁ"
)

var (
	isUserIn            map[string]bool
	guildChanToGameData map[string]*wfGame.Game

	env map[string]string
	emj map[string]string
	rg  []wfGame.RoleGuide
)

func readJSON() {
	readCommandJSON()
	readRuleJSON()
	readNoteJSON()
	readBackgroundJSON()
	readVictoryConditionJSON()
}

func init() {
	env = EnvInit()
	emj = EmojiInit()
	RoleGuideInit(&rg)
	readJSON()

	isUserIn = make(map[string]bool)
	guildChanToGameData = make(map[string]*wfGame.Game)
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
	enterUserIDChan := make(chan string, 1)
	quitUserIDChan := make(chan string)
	gameStartedChan := make(chan bool)
	curGame := wfGame.NewGame(m.GuildID, m.ChannelID, m.Author.ID, s, rg, emj, enterUserIDChan, quitUserIDChan, gameStartedChan)
	// Mutex 필요할 것으로 예상됨.
	guildChanToGameData[m.GuildID+m.ChannelID] = curGame
	for {
		select {
		case curUID := <-curGame.EnterUserIDChan:
			isUserIn[curUID] = true
		case curUID := <-curGame.QuitUserIDChan:
			delete(isUserIn, curUID)
		case <-curGame.GameStartedChan:
			return
		}
	}
}

// messageCreate() 입력한 메시지를 처리하는 함수
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	// 명령어모음
	if printHelpList(s, m) {
		return
	}
	switch m.Content {
	case "ㅁ시작":
		if guildChanToGameData[m.GuildID+m.ChannelID] != nil {
			s.ChannelMessageSend(m.ChannelID, "게임을 진행중인 채널입니다.")
			return
		}
		if isUserIn[m.Author.ID] {
			s.ChannelMessageSend(m.ChannelID, "게임을 진행중인 사용자입니다.")
			return
		}
		isUserIn[m.Author.ID] = true
		go startgame(s, m)
	case "ㅁ강제종료":
		if isUserIn[m.Author.ID] {
			s.ChannelMessageSend(m.ChannelID, "3초 후 게임을 강제종료합니다.")
			time.Sleep(3 * time.Second)
			g := guildChanToGameData[m.GuildID+m.ChannelID]
			if m.Author.ID != g.MasterID {
				return
			}
			for _, user := range g.UserList {
				delete(isUserIn, user.UserID)
			}
			delete(guildChanToGameData, m.GuildID+m.ChannelID)
			g.CanFunc()
			s.ChannelMessageSend(m.ChannelID, "게임을 강제종료 했습니다.")
		}
	case "ㅁ투표":
		uidChan := make(chan string, 7)
		thisGame := wfGame.NewGame(m.GuildID, m.ChannelID, m.Author.ID, s, rg, emj, uidChan, nil, nil)
		thisGame.UserList = append(thisGame.UserList, wfGame.NewUser(m.Author.ID, "jae-kim", m.ChannelID, m.ChannelID))
		thisGame.UserList = append(thisGame.UserList, wfGame.NewUser(m.Author.ID, "juhur", m.ChannelID, m.ChannelID))
		thisGame.UserList = append(thisGame.UserList, wfGame.NewUser(m.Author.ID, "min-jo", m.ChannelID, m.ChannelID))
		thisGame.UserList = append(thisGame.UserList, wfGame.NewUser(m.Author.ID, "kalee", m.ChannelID, m.ChannelID))
		thisGame.UserList = append(thisGame.UserList, wfGame.NewUser(m.Author.ID, "apple", m.ChannelID, m.ChannelID))
		thisGame.UserList = append(thisGame.UserList, wfGame.NewUser(m.Author.ID, "banana", m.ChannelID, m.ChannelID))

		temp2 := &wfGame.StateBeforeVote{thisGame}
		temp2.InitState()

		voted_list := make([]int, len(thisGame.UserList))
		temp := &wfGame.StateVote{thisGame, voted_list, len(thisGame.UserList), 0}
		guildChanToGameData[m.GuildID+m.ChannelID] = thisGame
		isUserIn[m.Author.ID] = true
		thisGame.CurState = temp
		wfGame.VoteProcess(s, thisGame)
	case "!test":
		str := rg[3].RoleGuide[0]
		s.ChannelMessageSend(m.ChannelID, str)
	}
}

// messageReactionAdd 함수는 인게임 버튼 이모지 상호작용 처리를 위한 이벤트 핸들러 함수입니다.
func messageReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	//fmt.Println(r.UserID, r.MessageID, r.ChannelID, r.GuildID)

	// 봇 자기자신의 리액션 무시.
	if r.UserID == s.State.User.ID {
		return
	}
	// 게임 참가중이 아닌 사용자의 리액션 무시.
	// 단, 참가자가 아니면 참가 가능해야 함. 무시해버리면 참가 못 함.
	if !(isUserIn[r.UserID] || (!isUserIn[r.UserID] && r.Emoji.Name == emj["YES"])) {
		return
	}
	g := guildChanToGameData[r.GuildID+r.ChannelID]
	if g == nil {
		return
	}
	isUserIn[r.UserID] = true
	// 숫자 이모지 선택.
	for i := 1; i < 10; i++ {
		var ch rune
		ch = '0' + rune(i)
		emjID := "n" + string(ch)
		if r.Emoji.Name == emj[emjID] {
			go g.CurState.PressNumBtn(s, r, i)
			break
		}
	}
	switch r.Emoji.Name {
	case emj["DISCARD"]:
		// 쓰레기통 이모지 선택.
		g.CurState.PressDisBtn(s, r)
	case emj["YES"]:
		// O 이모지 선택.
		g.CurState.PressYesBtn(s, r)
	case emj["NO"]:
		// X 이모지 선택.
		g.CurState.PressNoBtn(s, r)
	case emj["LEFT"]:
		// 왼쪽 화살표 선택.
		g.CurState.PressDirBtn(s, r, -1)
	case emj["RIGHT"]:
		// 오른쪽 화살표 선택.
		g.CurState.PressDirBtn(s, r, 1)
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

// 게임진행에 관련된 명령어가 입력될시 각 명령어에 해당하는 메시지를 출력하는 함수
func printHelpList(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	// "ㅁ직업소개 <직업명>" or "ㅁ직업소개 모두"
	if printRoleInfo(s, m) {
		return true
	}
	switch m.Content {
	case prefix + "도움말":
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed(backgroundTitle, backgroundMsg))
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed(noteTitle, noteMsg))
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed(ruleTitle, ruleMsg))
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed(commandTitle, commandMsg))
	case prefix + "명령어":
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed(commandTitle, commandMsg))
	case prefix + "게임배경":
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed(backgroundTitle, backgroundMsg))
	case prefix + "게임방법":
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed(ruleTitle, ruleMsg))
	case prefix + "참고":
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed(noteTitle, noteMsg))
	case prefix + "승리조건":
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed(VCTitle, VCMsg))
	case prefix + "직업목록":
		printRoleList(s, m)
	case prefix + "능력순서":
		printSkillOrder(s, m)
	default:
		return false
	}
	return true
}
