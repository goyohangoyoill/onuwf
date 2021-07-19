/* onuwf 는 보드게임 "한밤의 늑대인간" 을 디스코드 봇으로 구현하는 프로젝트입니다. */

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	wfGame "onuwf.com/game"
	util "onuwf.com/util"

	"github.com/bwmarrin/discordgo"
)

const (
	prefix = "ㅁ"
)

var (
	isUserIn            map[string]bool
	uidToGameData       map[string]*wfGame.Game
	guildChanToGameData map[string]*wfGame.Game

	env map[string]string
	emj map[string]string
	rg  []wfGame.RoleGuide
)

/*
type LoadDBInfo struct {
	MatchedUserList []*wfGame.User
	LastRoleSeq     []wfGame.Role //User로

}

type SaveDBInfo struct {
	CurUserList []*wfGame.User
	CurRoleSeq  []wfGame.Role
	mUserID     string
}
*/
func init() {
	env = util.EnvInit()
	emj = util.EmojiInit()
	RoleGuideInit(&rg)
	util.ReadJSON(rg, prefix)
	//util.MongoConn(env)

	isUserIn = make(map[string]bool)
	guildChanToGameData = make(map[string]*wfGame.Game)
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
	enterUserIDChan := make(chan string, 1)
	quitUserIDChan := make(chan string)
	gameStartedChan := make(chan bool)
	curGame := wfGame.NewGame(m.GuildID, m.ChannelID, m.Author.ID, s, rg, emj, enterUserIDChan, quitUserIDChan, gameStartedChan)
	// Mutex 필요할 것으로 예상됨.
	guildChanToGameData[m.GuildID+m.ChannelID] = curGame
	uidToGameData[m.Author.ID] = curGame
	LoadEnterUser(curGame, m.Author.ID)
	for {
		select {
		case curUID := <-curGame.EnterUserIDChan:
			isUserIn[curUID] = true
			guildChanToGameData[m.GuildID+curUID] = curGame
			uidToGameData[curUID] = curGame
			LoadEnterUser(curGame, curUID)
		case curUID := <-curGame.QuitUserIDChan:
			delete(isUserIn, curUID)
			delete(uidToGameData, curUID)
		case <-curGame.GameStartedChan:
			//LoadDB(curGame)
			SaveStartDB(curGame)
			return
		}
	}
}

func LoadEnterUser(g *wfGame.Game, uid string) {
	conn, ctx := util.MongoConn(env)
	m := false
	if g.MasterID == uid {
		m = true
	}
	lUser, p := util.LoadEachUser(uid, m, "User", conn.Database("ONUWF"), ctx)
	cUser := g.FindUserByUID(uid)
	if p == true {
		wfGame.UpdateUser(cUser, lUser.Nick, lUser.Title)
		//g.DelUserByID(uid)
		//g.UserList = append(g.UserList, cUser)
		if m == true {
			g.FormerRole = lUser.LastRole
			/*
				lenRole := len(lUser.LastRole)
				for i := 0; i < lenRole; i++ {
					g.AddRole(lUser.LastRole[i])
				}
			*/
		}
	}
	fmt.Println(g.UserList[0].Nick(), g.UserList[0].Title())
	fmt.Println(g.FormerRole)
}

/*
func LoadDB(g *wfGame.Game) {
	conn, ctx := util.MongoConn(env)
	uLen := len(g.UserList)
	sDB := util.SaveDBInfo{g.UserList, RoleID, g.MasterID}
	lDB := util.LoadUser(sDB, "User", conn.Database("ONUWF"), ctx)
	mLen := len(lDB.LoadedUser)
	for i := 0; i < uLen; i++ {
		//	a := util.CreateUser(g.UserList[i], "User", conn.Database("ONUWF"), ctx)
		for j := 0; j < mLen; j++ {

		}
	}

}
*/
func SaveStartDB(g *wfGame.Game) {
	conn, ctx := util.MongoConn(env)
	rLen := len(g.RoleView)
	RoleID := make([]int, rLen)
	for i := 0; i < rLen; i++ {
		RoleID[i] = g.RoleView[i].ID()
	}
	sDB := util.SaveDBInfo{g.UserList, RoleID, g.MasterID}

	a := util.SetStartUser(sDB, "User", conn.Database("ONUWF"), ctx)
	/*
		rLen := len(g.RoleView)
		RoleID := make([]int, rLen)
		for i := 0; i < rLen; i++ {
			RoleID[i] = g.RoleView[i].ID()
		}
		ret := util.SaveRole(&RoleID, "Game", conn.Database("ONUWF"), ctx)
	*/
	fmt.Println(a)
}

// messageCreate() 입력한 메시지를 처리하는 함수
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	// 명령어모음
	if util.PrintHelpList(s, m, rg, prefix) {
		return
	}
	switch m.Content {
	case prefix + "시작":
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
	case prefix + "강제종료":
		if isUserIn[m.Author.ID] {
			s.ChannelMessageSend(m.ChannelID, "3초 후 게임을 강제종료합니다.")
			time.Sleep(3 * time.Second)
			g := guildChanToGameData[m.GuildID+m.ChannelID]
			if m.Author.ID != g.MasterID {
				return
			}
			for _, user := range g.UserList {
				delete(isUserIn, user.UserID)
				delete(uidToGameData, user.UserID)
			}
			delete(guildChanToGameData, m.GuildID+m.ChannelID)
			g.CanFunc()
			s.ChannelMessageSend(m.ChannelID, "게임을 강제종료 했습니다.")
		}
	case prefix + "확인":
		g := guildChanToGameData[m.GuildID+m.ChannelID]
		if g != nil {
			Server, _ := s.State.Guild(m.GuildID)
			Channel, _ := s.State.Channel(m.ChannelID)
			msg := "----------------------------------------------------\n"
			msg += "> 현재 서버: " + Server.Name + "\n"
			msg += "> 현재 채널: " + Channel.Name + "\n"
			msg += "> 현재 유저 수: " + strconv.Itoa(len(g.UserList)) + "\n"
			msg += "----------------------------------------------------\n"
			for i, user := range g.UserList {
				msg += "< " + strconv.Itoa(i+1) + "번 유저 `" + user.Nick() + "` >\n"
				msg += "원래직업: " + g.GetOriRole(user.UserID).String() + "\n"
				msg += "현재직업: " + g.GetRole(user.UserID).String() + "\n"
			}
			msg += "< 버려진 직업들 >\n"
			for i := 0; i < 3; i++ {
				msg += g.GetDisRole(i).String() + " "
			}
			msg += "\n"
			msg += "----------------------------------------------------\n"
			msg += "로그 메시지 :\n"
			for _, text := range g.LogMsg {
				msg += text + "\n"
			}
			msg += "----------------------------------------------------\n"
			s.ChannelMessageSend(m.ChannelID, msg)
		}
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
	g := uidToGameData[r.UserID]
	if g == nil {
		g = guildChanToGameData[r.GuildID+r.ChannelID]
		if g == nil {
			return
		}
	}
	isUserIn[r.UserID] = true
	// 숫자 이모지 선택.
	for i := 1; i < 10; i++ {
		emjID := "n" + strconv.Itoa(i)
		if r.Emoji.Name == emj[emjID] {
			go g.CurState.PressNumBtn(s, r.MessageReaction, i)
			break
		}
	}
	switch r.Emoji.Name {
	case emj["DISCARD"]:
		// 쓰레기통 이모지 선택.
		g.CurState.PressDisBtn(s, r.MessageReaction)
	case emj["YES"]:
		// O 이모지 선택.
		g.CurState.PressYesBtn(s, r.MessageReaction)
	case emj["NO"]:
		// X 이모지 선택.
		g.CurState.PressNoBtn(s, r.MessageReaction)
	case emj["LEFT"]:
		// 왼쪽 화살표 선택.
		g.CurState.PressDirBtn(s, r.MessageReaction, -1)
	case emj["RIGHT"]:
		// 오른쪽 화살표 선택.
		g.CurState.PressDirBtn(s, r.MessageReaction, 1)
	case emj["BOOKMARK"]:
		g.CurState.PressBmkBtn(s, r.MessageReaction)
	}
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
