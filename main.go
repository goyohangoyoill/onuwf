/* onuwf 는 보드게임 "한밤의 늑대인간" 을 디스코드 봇으로 구현하는 프로젝트입니다. */

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	embed "github.com/clinet/discordgo-embed"
	wfGame "github.com/goyohangoyoill/onuwf/game"
	"github.com/goyohangoyoill/onuwf/util"
	"github.com/goyohangoyoill/onuwf/util/json"

	"github.com/bwmarrin/discordgo"
)

var (
	isUserIn            map[string]bool
	uidToGameData       map[string]*wfGame.Game
	guildChanToGameData map[string]*wfGame.Game
	fqChanMap           map[string]chan bool
	env                 map[string]string
	emj                 map[string]string
	rg                  []json.RoleGuide
	config              json.Config
	isNickChange        map[string]bool
	chNick              map[string]chan string
	globalStatus        string
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "시작",
			Description: "한밤의 늑대인간 게임 시작",
		},
		{
			Name:        "강제종료",
			Description: "현재 채널에서 실행중인 게임 강제종료",
		},
		{
			Name:        "관전",
			Description: "현재 진행되고 있는 게임 정보를 DM으로 불러옵니다.",
		},
		{
			Name:        "내정보",
			Description: "내 프로필 정보를 불러옵니다.",
		},
		{
			Name:        "도움말",
			Description: "도움말 불러오기",
		},
		{
			Name:        "명령어",
			Description: "명령어 불러오기",
		},
		{
			Name:        "help",
			Description: "명령어 불러오기",
		},
		{
			Name:        "게임배경",
			Description: "게임배경 불러오기",
		},
		{
			Name:        "게임방법",
			Description: "게임방법 불러오기",
		},
		{
			Name:        "참고",
			Description: "참고 불러오기",
		},
		{
			Name:        "승리조건",
			Description: "승리조건 불러오기",
		},
		{
			Name:        "직업목록",
			Description: "직업목록 불러오기",
		},
		{
			Name:        "직업순서",
			Description: "직업순서 불러오기",
		},
		{
			Name:        "직업서순",
			Description: "직업순서 역순으로 불러오기",
		},
		{
			Name:        "나무위키",
			Description: "나무위키 링크 불러오기",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"시작":   startGameHandler,
		"강제종료": forceStopGameHandler,
		"관전":   showGameStateHandler,
		"내정보":  myInfoHandler,
		"도움말":  helpHandler,
		"명령어":  helpHandler,
		"help": helpHandler,
		"게임배경": helpHandler,
		"게임방법": helpHandler,
		"참고":   helpHandler,
		"승리조건": helpHandler,
		"직업순서": helpHandler,
		"직업서순": helpHandler,
		"나무위키": helpHandler,
	}
)

func helpHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {

	if json.PrintHelpList(s, i, rg, "") {
		return
	}
}

func init() {
	env = json.EnvInit()
	emj = json.EmojiInit()
	json.RoleGuideInit(&rg)
	config = json.ReadConfigJson()
	json.ReadJSON(rg, config.Prefix)

	isUserIn = make(map[string]bool)
	guildChanToGameData = make(map[string]*wfGame.Game)
	uidToGameData = make(map[string]*wfGame.Game)
	fqChanMap = make(map[string]chan bool)
	isNickChange = make(map[string]bool)
	chNick = make(map[string]chan string)
	globalStatus = "!도움말 !명령어"
}

func main() {
	dg, err := discordgo.New("Bot " + env["dgToken"])
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
	dg.AddHandler(messageCreate)
	dg.AddHandler(messageReactionAdd)
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := dg.ApplicationCommandCreate(dg.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	log.Println("Removing commands...")
	// // We need to fetch the commands, since deleting requires the command ID.
	// // We are doing this from the returned commands on line 375, because using
	// // this will delete all the commands, which might not be desirable, so we
	// // are deleting only the commands that we added.
	// registeredCommands, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
	// if err != nil {
	// 	log.Fatalf("Could not fetch registered commands: %v", err)
	// }

	for _, v := range registeredCommands {
		err := dg.ApplicationCommandDelete(dg.State.User.ID, "", v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}

	log.Println("Gracefully shutting down.")
	_ = dg.Close()
}

func startgame(s *discordgo.Session, m *discordgo.InteractionCreate, isTest bool) {
	enterUserIDChan := make(chan string, 1)
	quitUserIDChan := make(chan string)
	gameStartedChan := make(chan bool)
	fqChanMap[m.GuildID+m.ChannelID] = make(chan bool, 1)
	curGame := wfGame.NewGame(m.GuildID, m.ChannelID, m.User.ID, s, rg, emj, config, enterUserIDChan, quitUserIDChan, gameStartedChan, env, isTest)
	// Mutex 필요할 것으로 예상됨.
	guildChanToGameData[m.GuildID+m.ChannelID] = curGame
	uidToGameData[m.User.ID] = curGame
	flag := false
	// juhur comment out
	for {
		if flag {
			break
		}
		select {
		case curUID := <-curGame.EnterUserIDChan:
			isUserIn[curUID] = true
			guildChanToGameData[m.GuildID+curUID] = curGame
			uidToGameData[curUID] = curGame
			// juhur comment out
		case curUID := <-curGame.QuitUserIDChan:
			delete(isUserIn, curUID)
			delete(uidToGameData, curUID)
		case _ = <-curGame.GameStartedChan:
			flag = true
			// juhur comment out
			SaveStartDB(curGame)
		}
	}
	<-curGame.GameStartedChan
	fqChanMap[m.GuildID+m.ChannelID] <- true
	g := guildChanToGameData[m.GuildID+m.ChannelID]
	if g == nil {
		<-fqChanMap[m.GuildID+m.ChannelID]
		return
	}
	// 여기에 DB 갱신 넣으면 됨.
	SaveEndDB(curGame)
	for _, user := range g.UserList {
		delete(isUserIn, user.UserID)
		delete(uidToGameData, user.UserID)
	}
	delete(guildChanToGameData, m.GuildID+m.ChannelID)
	g.CanFunc()
	s.ChannelMessageSend(m.ChannelID, "게임이 종료 되었습니다.")
	<-fqChanMap[m.GuildID+m.ChannelID]
}

// SaveStartDB : 게임 시작 시 save (user nick, lastrole 정보 저장)
func SaveStartDB(g *wfGame.Game) {
	conn, ctx := util.MongoConn(env)
	rLen := len(g.RoleView)
	RoleID := make([]int, rLen)
	// 게임 시작 시 설정 직업 정보를 가진 배열 초기
	for i := 0; i < rLen; i++ {
		RoleID[i] = g.RoleView[i].ID()
	}
	UserInfo := make([]*util.UserData, 0)
	uLen := len(g.UserList)
	for i := 0; i < uLen; i++ {
		UserInfo = append(UserInfo, &util.UserData{UID: g.UserList[i].UserID, Nick: g.UserList[i].Nick()})
	}
	sDB := util.SaveDBInfo{CurUserList: UserInfo, CurRoleSeq: RoleID, MUserID: g.MasterID}
	util.SetStartUser(sDB, "User", conn.Database("ONUWF"), ctx)
}

func SaveUserInit(g *wfGame.Game) []util.User {
	uLen := len(g.UserList)
	users := make([]util.User, 0, uLen)
	for i := 0; i < uLen; i++ {
		users = saveUser(g, i, false, users)
	}
	return users
}

func saveUser(g *wfGame.Game, i int, win bool, users []util.User) []util.User {
	user := util.User{}
	user.UID = g.UserList[i].UserID
	user.Nick = g.UserList[i].Nick()
	user.OriRole = g.GetOriRole(g.UserList[i].UserID).String()
	user.LastRole = g.GetRole(g.UserList[i].UserID).String()
	if (g.GetRole(g.UserList[i].UserID).String() == (&wfGame.Werewolf{}).String()) || (g.GetRole(g.UserList[i].UserID).String() == (&wfGame.Minion{}).String()) {
		win = g.WerewolfTeamWin
	} else if (g.GetRole(g.UserList[i].UserID).String()) == (&wfGame.Tanner{}).String() {
		win = g.TannerTeamWin
	} else {
		win = g.VillagerTeamWin
	}
	user.IsWin = win
	users = append(users, user)
	return users
}

func SaveGameInit(g *wfGame.Game) util.GameData {
	sGame := util.GameData{}
	sGame.GuildID = g.GuildID
	sGame.ChanID = g.ChanID
	sGame.MasterID = g.MasterID
	RoleList := make([]string, 0, len(g.RoleView))
	for i := 0; i < len(g.RoleView); i++ {
		RoleList = append(RoleList, g.RoleView[i].String())
	}
	sGame.RoleList = RoleList
	sGame.UserList = SaveUserInit(g)
	disRole := make([]string, 0, len(g.DisRole))
	oriDisRole := make([]string, 0, len(g.OriDisRole))
	for i := 0; i < len(g.DisRole); i++ {
		disRole = append(disRole, g.DisRole[i].String())
		oriDisRole = append(oriDisRole, g.OriDisRole[i].String())
	}
	sGame.OriDisRole = oriDisRole
	sGame.LastDisRole = disRole

	return sGame
}

func SaveEndDB(g *wfGame.Game) {
	conn, ctx := util.MongoConn(env)
	sGame := SaveGameInit(g)
	t := time.Now()
	curGameOID := util.SaveGame(sGame, t, "Game", conn.Database("ONUWF"), ctx)
	uLen := len(g.UserList)
	for i := 0; i < uLen; i++ {
		win := false
		mostVoted := false
		if (g.GetRole(g.UserList[i].UserID).String() == (&wfGame.Werewolf{}).String()) || (g.GetRole(g.UserList[i].UserID).String() == (&wfGame.Minion{}).String()) {
			win = g.WerewolfTeamWin
		} else if (g.GetRole(g.UserList[i].UserID).String()) == (&wfGame.Tanner{}).String() {
			win = g.TannerTeamWin
		} else {
			win = g.VillagerTeamWin
		}
		if g.MostVoted != nil {
			if g.UserList[i].UserID == g.MostVoted.UserID {
				mostVoted = true
			}
		}
		lUser, _ := util.LoadEachUser(g.UserList[i].UserID, true, "User", conn.Database("ONUWF"), ctx)
		util.SaveEachUser(&lUser, curGameOID, win, mostVoted, t, "User", conn.Database("ONUWF"), ctx)
	}
}

// messageCreate() 입력한 메시지를 처리하는 함수
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	// 명령어모음
	if isNickChange[m.Author.ID] {
		chNick[m.Author.ID] <- m.Content
		return
	}
	switch m.Content {
	case config.Prefix + "확인":
		g := guildChanToGameData[m.GuildID+m.ChannelID]
		if g == nil {
			return
		}
		if !g.IsTest {
			return
		}
		if len(g.OriRoleIdxTable) == 0 {
			return
		}
		Server, _ := s.State.Guild(m.GuildID)
		Channel, _ := s.State.Channel(m.ChannelID)
		msg := "----------------------------------------------------\n"
		msg += "> 현재 서버: " + Server.Name + "\n"
		msg += "> 현재 채널: " + Channel.Name + "\n"
		msg += "> 현재 유저 수: " + strconv.Itoa(len(g.UserList)) + "\n"
		msg += "----------------------------------------------------\n"
		s.ChannelMessageSend(m.ChannelID, msg)
		g.SendLogMsg(m.ChannelID)
	case config.Prefix + "닉네임":
		isNickChange[m.Author.ID] = true
		chNick[m.Author.ID] = make(chan string)
		chTimeout := make(chan bool)
		dmChan, _ := s.UserChannelCreate(m.Author.ID)
		msg := "닉네임을 변경하려면 " + strconv.Itoa(config.NickChangeSec) + "초 안에 입력해주세요."
		s.ChannelMessageSend(dmChan.ID, msg)
		go func(chan bool) {
			time.Sleep(time.Duration(config.NickChangeSec) * time.Second)
			chTimeout <- true
		}(chTimeout)
		select {
		case nick := <-chNick[m.Author.ID]:
			conn, mgctx := util.MongoConn(env)
			user, _ := util.LoadEachUser(m.Author.ID, false, "User", conn.Database("ONUWF"), mgctx)
			util.SetUserNick(&user, nick, conn.Database("ONUWF"), mgctx)
			s.ChannelMessageSend(dmChan.ID, "닉네임을 "+nick+"으로 변경했습니다.")
			delete(chNick, m.Author.ID)
			isNickChange[m.Author.ID] = false
		case _ = <-chTimeout:
			s.ChannelMessageSend(dmChan.ID, "닉네임을 변경하지 않았습니다.")
			delete(chNick, m.Author.ID)
			isNickChange[m.Author.ID] = false
		}
	}
}

func myInfoHandler(s *discordgo.Session, m *discordgo.InteractionCreate) {
	conn, mgctx := util.MongoConn(env)
	user, _ := util.LoadEachUser(m.User.ID, false, "User", conn.Database("ONUWF"), mgctx)
	if user.CntPlay == 0 {
		return
	}
	myInfoEmbed := embed.NewEmbed()
	myInfoEmbed.SetTitle("한밤의 늑대인간 유저정보")
	myInfoEmbed.AddField("닉네임", user.Nick)
	if len(user.Title) > 0 {
		myInfoEmbed.AddField("칭호", user.Title)
	}
	myInfoEmbed.AddField("게임횟수", strconv.Itoa(user.CntPlay)+"회")
	myInfoEmbed.AddField("승리횟수", strconv.Itoa(user.CntWin)+"회(승률:"+strconv.Itoa(user.CntWin*100/user.CntPlay)+"%)")
	myInfoEmbed.AddField("최근게임시간", user.RecentGameTime.String())
	s.ChannelMessageSendEmbed(m.ChannelID, myInfoEmbed.MessageEmbed)
}

func showGameStateHandler(s *discordgo.Session, m *discordgo.InteractionCreate) {
	g := guildChanToGameData[m.GuildID+m.ChannelID]
	if g == nil {
		return
	}
	if len(g.OriRoleIdxTable) == 0 {
		return
	}
	if isUserIn[m.User.ID] {
		s.ChannelMessageSend(m.ChannelID, "게임에 참가중인 사람은 관전할 수 없습니다.")
		return
	}
	dmChan, _ := s.UserChannelCreate(m.User.ID)
	g.SendLogMsg(dmChan.ID)
	s.ChannelMessageSend(dmChan.ID, "진행상황을 더 알고싶으면 게임중인 채널에서 `!관전` 을 다시 입력하세요")
}

func forceStopGameHandler(s *discordgo.Session, m *discordgo.InteractionCreate) {
	if isUserIn[m.User.ID] {
		curChan := fqChanMap[m.GuildID+m.ChannelID]
		// Mutex Lock
		curChan <- true
		g := guildChanToGameData[m.GuildID+m.ChannelID]
		if g == nil {
			<-curChan
			return
		}
		if m.User.ID != g.MasterID {
			<-curChan
			return
		}
		s.ChannelMessageSend(m.ChannelID, "3초 후 게임을 강제종료합니다.")
		time.Sleep(3 * time.Second)
		g = guildChanToGameData[m.GuildID+m.ChannelID]
		if g == nil {
			<-curChan
		}
		for _, user := range g.UserList {
			delete(isUserIn, user.UserID)
			delete(uidToGameData, user.UserID)
		}
		delete(guildChanToGameData, m.GuildID+m.ChannelID)
		g.CanFunc()
		s.ChannelMessageSend(m.ChannelID, "게임을 강제종료 했습니다.")
		// Mutex Release
		<-curChan
	}
}

func startGameHandler(s *discordgo.Session, m *discordgo.InteractionCreate) {
	if guildChanToGameData[m.GuildID+m.ChannelID] != nil {
		s.ChannelMessageSend(m.ChannelID, "게임을 진행중인 채널입니다.")
		return
	}
	if isUserIn[m.User.ID] {
		s.ChannelMessageSend(m.ChannelID, "게임을 진행중인 사용자입니다.")
		return
	}
	isUserIn[m.User.ID] = true
	go startgame(s, m, false)
}

// messageReactionAdd 함수는 인게임 버튼 이모지 상호작용 처리를 위한 이벤트 핸들러 함수입니다.
func messageReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
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
	for i := 1; i < 10; i++ {
		emjID := "n" + strconv.Itoa(i)
		if r.Emoji.Name == emj[emjID] {
			g.CurState.PressNumBtn(s, r.MessageReaction, i)
			break
		}
	}
	switch r.Emoji.Name {
	case emj["DISCARD"]:
		// 🚮
		g.CurState.PressDisBtn(s, r.MessageReaction)
	case emj["YES"]:
		// ⭕️
		g.CurState.PressYesBtn(s, r.MessageReaction)
	case emj["NO"]:
		// ❌
		g.CurState.PressNoBtn(s, r.MessageReaction)
	case emj["LEFT"]:
		// ◀️
		g.CurState.PressDirBtn(s, r.MessageReaction, -1)
	case emj["RIGHT"]:
		// ▶️
		g.CurState.PressDirBtn(s, r.MessageReaction, 1)
	case emj["BOOKMARK"]:
		g.CurState.PressBmkBtn(s, r.MessageReaction)
	}
	s.UpdateListeningStatus(globalStatus)
}
