package game

import (
	"context"
	"fmt"
	"sort"

	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
	util "github.com/goyohangoyoill/onuwf/util"
	json "github.com/goyohangoyoill/onuwf/util/json"
)

// Game 구조체는 게임 진행을 위한 정보를 담고 있는 스트럭처
type Game struct {
	// 게임이 테스트중인지 체크하는 변수
	IsTest bool
	// 게임을 강제 종료하기 위한 컨텍스트.
	Ctx context.Context
	// 게임을 강제 종료하기 위한 캔슬함수.
	CanFunc context.CancelFunc
	// 현재 게임이 진행중인 서버의 GID
	GuildID string
	// 현재 게임이 진행중인 채널의 CID
	ChanID string
	// 게임을 생성한 방장의 UID
	MasterID string
	// 게임 상태 표시 메시지의 mid
	GameStateMID string
	// 현재 게임의 진행시점
	CurState State
	// 현재 게임의 참가자들
	UserList []*User
	// 현재 게임에서 role_guide.json 순서대로(role ID 순서대로) 추가, 중복제거 된 직업들의 목록
	RoleSeq []Role
	// 현재 게임에서 사용중인 사용자에게 보여줄 중복 직업들의 목록, 정렬 안됨
	RoleView []Role
	// Role을 User별로 매핑시킨 인덱스 테이블
	// RoleSeq 사용
	// <usage : roleIdxTable[userIdx][roleIdx]>
	roleIdxTable    [][]int
	OriRoleIdxTable [][]int
	// 게임에서 버려진 직업 목록
	DisRole    []Role
	OriDisRole []Role
	// 게임에서 사용하는 세션
	Session *discordgo.Session
	// 게임 진행 상황을 기록하는 로그 메시지 배열
	LogMsg []string
	// 이모지 정보
	Emj map[string]string
	// 환경설정 정보
	config json.Config
	// 직업의 대한 소개 및 정보
	RG []json.RoleGuide
	// 유저 입장, 퇴장 시 ID가 전달되는 채널
	EnterUserIDChan, QuitUserIDChan chan string
	// 게임이 시작되면 신호가 전달되는 채널
	GameStartedChan chan bool

	FormerRole []int
	// 마을주민팀 승패여부
	VillagerTeamWin bool
	// 늑대인간팀 승패여부
	WerewolfTeamWin bool
	// 무두장이 승패여부
	TannerTeamWin bool
	// 게임에서 db에 접근해야할 경우 사용하는 환경변수
	env map[string]string

	MostVoted *User
}

// NewGame : Game 스트럭처를 생성하는 생성자,
func NewGame(gid, cid, muid string, s *discordgo.Session, rg []json.RoleGuide, emj map[string]string, config json.Config, enterUserIDChan, quitUserIDChan chan string, gameStartedChan chan bool, env map[string]string, isTest bool) (g *Game) {
	g = &Game{}
	g.GuildID = gid
	g.ChanID = cid
	g.MasterID = muid
	g.Session = s
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	g.Ctx = ctx
	g.CanFunc = cancel
	g.RG = rg
	g.Emj = emj
	g.env = env
	g.config = config
	g.EnterUserIDChan = enterUserIDChan
	g.QuitUserIDChan = quitUserIDChan
	g.GameStartedChan = gameStartedChan
	var maxrole int
	for _, roleItem := range rg {
		maxrole += roleItem.Max
	}
	g.UserList = make([]*User, 0, maxrole-3)
	g.RoleSeq = make([]Role, 0, len(rg))
	g.RoleView = make([]Role, 0, maxrole)
	g.DisRole = make([]Role, 3)
	g.OriDisRole = make([]Role, 3)
	g.LogMsg = make([]string, 0)
	g.SetUserByID(muid)
	g.CurState = &Prepare{g, 0, nil, nil, false}
	g.GameStateMID = ""
	g.IsTest = isTest
	g.CurState.InitState()
	return
}

// SendVoteMsg 는 현재 참가자 모두에게 DM으로 투표 용지를 전송하고,
// 각각의 투표 용지별로 UserList index 순서에 맞춰 MsgID 배열을 반환해주는 함수이다.
func (g *Game) SendVoteMsg(s *discordgo.Session) (messageIDs []string) {
	messageIDs = make([]string, len(g.UserList))
	for i, me := range g.UserList {
		msg := ""
		userListExceptMe := append(g.UserList[:i], g.UserList[i:]...)
		for i := 0; i < 9; i++ {
			if i >= len(userListExceptMe) {
				break
			}
			msg += "`" + userListExceptMe[i].nick + "`\n"
		}
		mObj, _ := s.ChannelMessageSendEmbed(me.dmChanID, embed.NewGenericEmbed("투표 시작!", msg))
		messageIDs[i] = mObj.ID
	}
	return messageIDs
}

// IsDoppel 은 UserID 로 해당 유저가 도플갱어인지 확인하는 메소드입니다.
func (g *Game) IsDoppel(uid string) (res bool) {
	res = false
	uIdx := FindUserIdx(uid, g.UserList)
	for i := 0; i < len(g.RoleSeq); i++ {
		if g.OriRoleIdxTable[uIdx][i] == 2 {
			res = true
			break
		}
	}
	return res
}

// IsProtected 는 센티넬에 의해 보호받는 상태인지 확인하는 메소드입니다.
func (g *Game) IsProtected(uid string) bool {
	for _, user := range g.UserList {
		if uid == user.UserID && user.protected {
			return true
		}
	}
	return false
}

// SetUserByID 는 게임에 입장한 유저의 정보를 게임 데이터에 추가하는 함수입니다.
func (g *Game) SetUserByID(uid string) {
	if i := FindUserIdx(uid, g.UserList); i != -1 {
		return
	}
	if len(g.UserList) >= 10 {
		return
	}
	newOne := &User{}
	newOne.UserID = uid
	dgUser, _ := g.Session.User(uid)
	newOne.nick = dgUser.Username
	newOne.chanID = g.ChanID
	uChan, _ := g.Session.UserChannelCreate(uid)
	newOne.dmChanID = uChan.ID

	conn, ctx := util.MongoConn(g.env)
	// m은 master_user 식별 bool
	m := false
	if g.MasterID == uid {
		m = true
	}
	// p는 database에 존재여부 식별 bool
	lUser, p := util.LoadEachUser(uid, m, "User", conn.Database("ONUWF"), ctx)
	if p {
		newOne.nick, newOne.title = lUser.Nick, lUser.Title
		if m {
			g.FormerRole = lUser.LastRoleList
		}
	}

	g.EnterUserIDChan <- uid
	g.UserList = append(g.UserList, newOne)

}

// DelUserByID 는 입장되어 있는 유저의 정보를 모두 삭제해주는 함수입니다.
func (g *Game) DelUserByID(uid string) {
	if uid == g.MasterID {
		return
	}
	i := FindUserIdx(uid, g.UserList)
	if i == -1 {
		return
	}
	g.QuitUserIDChan <- uid
	g.UserList = append(g.UserList[:i], g.UserList[i+1:]...)
}

// RoleCount 함수는 직업을 가진 유저가 아닌 직업 자체의 갯수를 셈
func (g *Game) RoleCount(roleToFind Role, roleList []Role) int {
	cnt := 0
	for _, tmpRole := range roleList {
		if roleToFind.String() == tmpRole.String() {
			cnt++
		}
	}
	return cnt
}

// AddRole 함수는 RG에 사용할 roleindex 위치 값을 받아 RoleView와 RoleSeq에 role을 추가
func (g *Game) AddRole(roleIndex int) {
	// roleFactory에서 현재 roleindex 위치 값을 받아 role 생성
	roleToAdd := GenerateRole(roleIndex)
	// max 넘기면 초기화
	if g.RoleCount(roleToAdd, g.RoleView) == g.RG[roleIndex].Max {
		for i := 0; i < g.RG[roleIndex].Max; i++ {
			g.DelRole(roleIndex)
		}
		return
	}
	// RoleView에 추가된 role 개수가 max보다 작을 때만 추가
	g.RoleView = append(g.RoleView, roleToAdd)
	if FindRoleIdx(roleToAdd, g.RoleSeq) == -1 {
		g.RoleSeq = append(g.RoleSeq, roleToAdd)
		sort.Slice(g.RoleSeq, func(i, j int) bool {
			return g.RoleSeq[i].ID() < g.RoleSeq[j].ID()
		})
	}
}

// DelRole 함수는 RG에 사용할 roleindex 위치 값을 받아 RoleView와 RoleSeq에서 role을 삭제
func (g *Game) DelRole(roleIndex int) {
	// roleFactory에서 현재 roleindex 위치 값을 받아 role 생성
	roleToRemove := GenerateRole(roleIndex)
	// RoleView는 ununique sorted니까 첫번째로 나오는거 찾아서 지우기
	if i := FindRoleIdx(roleToRemove, g.RoleView); i != -1 {
		/* i + 1 이 인덱스 범위를 벗어나는가에 대한 고민이 필요함
		if i+1 == len(g.RoleView) {
			g.RoleView = g.RoleView[:i]
		} else {
			g.RoleView = append(g.RoleView[:i], g.RoleView[i+1:]...)
		}
		*/
		g.RoleView = append(g.RoleView[:i], g.RoleView[i+1:]...)
	}
	// RoleSeq는 unique unsorted니까 방금 지운 RoleView에 없으면 지우기
	if g.RoleCount(roleToRemove, g.RoleView) == 0 {
		if i := FindRoleIdx(roleToRemove, g.RoleSeq); i != -1 {
			/* 위 주석과 비슷한 문제로 범위 벗어나는지 확인해야함
			if i+1 == len(g.RoleView) {
				g.RoleSeq = g.RoleSeq[:i]
			} else {
				g.RoleSeq = append(g.RoleSeq[:i], g.RoleSeq[i+1:]...)
			}
			*/
			g.RoleSeq = append(g.RoleSeq[:i], g.RoleSeq[i+1:]...)
		}
	}
}

// FindUserByUID UID 로 user 인스턴스를 구하는 함수
func (g *Game) FindUserByUID(uid string) (target *User) {
	for i, item := range g.UserList {
		if item.UserID == uid {
			return g.UserList[i]
		}
	}
	return nil
}

// AppendLog 게임 로그에 메시지를 쌓는 함수.
func (g *Game) AppendLog(msg string) {
	if g.LogMsg == nil {
		g.LogMsg = make([]string, 0)
	}
	g.LogMsg = append(g.LogMsg, msg)
}

// GetRole 유저의 직업을 반환
func (g *Game) GetRole(uid string) Role {
	idx := FindUserIdx(uid, g.UserList)
	if idx == -1 {
		fmt.Println("user idx is -1")
		return nil
	}
	for i := 0; i < len(g.RoleSeq); i++ {
		if g.roleIdxTable[idx][i] > 0 {
			return g.RoleSeq[i]
		}
	}
	return nil
}

// GetOriRole 유저의 원래 직업을 반환
// 원래 직업이 도플갱어였다면 값이 2
func (g *Game) GetOriRole(uid string) Role {
	idx := FindUserIdx(uid, g.UserList)
	if idx == -1 {
		fmt.Println("user idx is -1")
		return nil
	}
	for i := 0; i < len(g.RoleSeq); i++ {
		if g.OriRoleIdxTable[idx][i] > 0 {
			if g.OriRoleIdxTable[idx][i] == 2 {
				return &Doppelganger{}
			}
			return g.RoleSeq[i]
		}
	}
	return nil
}

// 유저의 직업을 업데이트
func (g *Game) setRole(uid string, item Role) {
	userIdx := FindUserIdx(uid, g.UserList)
	roleIdx := FindRoleIdx(item, g.RoleSeq)
	if userIdx == -1 || roleIdx == -1 {
		fmt.Println("table idx is -1")
		return
	}
	loop := len(g.RoleSeq)
	for i := 0; i < loop; i++ {
		g.roleIdxTable[userIdx][i] = 0
	}
	g.roleIdxTable[userIdx][roleIdx] = 1
}

// 도플갱어인 유저의 직업을 업데이트
func (g *Game) setDplRole(uid string, item Role) {
	userIdx := FindUserIdx(uid, g.UserList)
	roleIdx := FindRoleIdx(item, g.RoleSeq)
	if userIdx == -1 || roleIdx == -1 {
		fmt.Println("table idx is -1")
		return
	}
	loop := len(g.RoleSeq)
	for i := 0; i < loop; i++ {
		g.OriRoleIdxTable[userIdx][i] = 0
		g.roleIdxTable[userIdx][i] = 0
	}
	g.OriRoleIdxTable[userIdx][roleIdx] = 2
	g.roleIdxTable[userIdx][roleIdx] = 1
}

// SetDisRole 버려진 직업을 업데이트
func (g *Game) SetDisRole(disRoleIdx int, item Role) {
	g.DisRole[disRoleIdx] = item
}

// SwapRoleFromUser 두 유저의 직업을 서로 교환
func (g *Game) SwapRoleFromUser(uid1, uid2 string) {
	role1 := g.GetRole(uid1)
	role2 := g.GetRole(uid2)
	g.setRole(uid1, role2)
	g.setRole(uid2, role1)
}

// GetDisRole : 버려진 직업 중 하나 확인.
func (g *Game) GetDisRole(disRoleIdx int) Role {
	return g.DisRole[disRoleIdx]
}

// SwapRoleFromDiscard 유저 직업과 버려진 직업을 교환.
func (g *Game) SwapRoleFromDiscard(uid string, disRoleIdx int) {
	role1 := g.DisRole[disRoleIdx]
	role2 := g.GetRole(uid)
	g.setRole(uid, role1)
	g.SetDisRole(disRoleIdx, role2)
}

// GetRoleUsers 특정 직업의 유저 목록 반환.
// ※ FindRoleIdx()에서 -1을 반환할 경우 panic 발생 (늑대인간을 제외한 직업을 사용할 때는 수정이 필요)
func (g *Game) GetRoleUsers(find Role) (result []*User) {
	result = make([]*User, 0)
	loop := len(g.UserList)
	idx := FindRoleIdx(find, g.RoleSeq)
	if idx == -1 {
		return result
	}
	for i := 0; i < loop; i++ {
		if g.roleIdxTable[i][idx] > 0 {
			result = append(result, g.UserList[i])
		}
	}
	return result
}

// GetOriRoleUsers 특정 원래 직업의 유저 목록 반환.
func (g *Game) GetOriRoleUsers(find Role) (result []*User) {
	result = make([]*User, 0)
	loop := len(g.UserList)
	idx := FindRoleIdx(find, g.RoleSeq)
	if idx == -1 {
		return result
	}
	for i := 0; i < loop; i++ {
		if g.OriRoleIdxTable[i][idx] > 0 {
			result = append(result, g.UserList[i])
		}
	}
	return result
}

// GetOriRoleUsersWithoutDpl 특정 원래 직업의 유저 목록 반환.
func (g *Game) GetOriRoleUsersWithoutDpl(find Role) (users []*User) {
	result := make([]*User, 0)
	loop := len(g.UserList)
	idx := FindRoleIdx(find, g.RoleSeq)
	for i := 0; i < loop; i++ {
		if g.OriRoleIdxTable[i][idx] == 1 {
			result = append(result, g.UserList[i])
		}
	}
	return result
}

// RotateAllUserRole  모든 사람들의 직업을 입장순서별로 한칸 회전.
// 문제 많아서 잠시 구현 보류
func (g *Game) RotateAllUserRole() {
	loop := len(g.UserList)
	shIdx := -1
	var shUID string
	for i, user := range g.UserList {
		if g.IsProtected(user.UserID) {
			shUID = user.UserID
			shIdx = i
		}
	}
	tmpRole := g.GetRole(g.UserList[loop-1].UserID)
	for i := 1; i < loop; i++ {
		srcUser := g.UserList[i-1]
		destUser := g.UserList[i]
		role := g.GetRole(srcUser.UserID)
		g.AppendLog("`" + srcUser.nick + "`의 직업 `" + role.String() + "` 은(는)\n다음 플레이어 `" + destUser.nick + "` 에게 주어졌습니다.")
		g.setRole(destUser.UserID, role)
	}
	if shIdx != -1 {
		g.SetProtect(shUID)
	}
	g.AppendLog("마지막 플레이어 `" + g.UserList[loop-1].nick + "`의 직업 `" + tmpRole.String() + "` 은(는)\n첫번째 플레이어 `" + g.UserList[0].nick + "` 에게 주어졌습니다.")
	g.setRole(g.UserList[0].UserID, tmpRole)
}

// SetProtect 유저에게 특수권한 부여
func (g *Game) SetProtect(uid string) {
	for _, user := range g.UserList {
		if uid == user.UserID {
			user.protected = true
			break
		}
	}
}

// CopyRole 특정 유저의 직업을 복사.
func (g *Game) CopyRole(destUID, srcUID string) {
	srcRole := g.GetRole(srcUID)
	g.setRole(destUID, srcRole)
}

// DplCopyRole 도플갱어인 유저의 직업을 다른 사람것으로 복사
func (g *Game) DplCopyRole(destUID, srcUID string) {
	srcRole := g.GetRole(srcUID)
	g.setDplRole(destUID, srcRole)
}

// FindUserIdx 유저의 인덱스 찾기를 위한 함수
func FindUserIdx(uid string, target []*User) int {
	for i, item := range target {
		if uid == item.UserID {
			return i
		}
	}
	return -1
}

// FindRoleIdx 직업의 인덱스 찾기를 위한 함수
func FindRoleIdx(r Role, target []Role) int {
	for i, item := range target {
		if r.String() == item.String() {
			return i
		}
	}
	return -1
}

// SendLogMsg 현재 게임의 로그 메시지를 전송하는 함수
func (g *Game) SendLogMsg(cid string) {
	s := g.Session
	tmpEmbed := embed.NewEmbed()
	tmpEmbed.SetTitle("직업 배정")
	roleListTitle := ""
	roleListMsg := ""
	for _, user := range g.UserList {
		roleListTitle = "`" + user.nick + "`"
		if g.IsProtected(user.UserID) {
			roleListTitle += " " + g.Emj["SHIELD"]
		}
		roleListMsg = "원래직업 : `" + g.GetOriRole(user.UserID).String() + "`"
		roleListMsg += g.getTeamMark(g.GetOriRole(user.UserID).String()) + "\n"
		roleListMsg += "현재직업 : `" + g.GetRole(user.UserID).String() + "`"
		roleListMsg += g.getTeamMark(g.GetRole(user.UserID).String())
		tmpEmbed.AddField(roleListTitle, roleListMsg)
	}
	tmpEmbed.InlineAllFields()
	disMsg := ""
	for i := 0; i < 3; i++ {
		disMsg += "`" + g.DisRole[i].String() + "`"
		disMsg += g.getTeamMark(g.DisRole[i].String()) + " "
	}
	tmpEmbed.AddField("버려진 직업들", disMsg)
	logMsg := ""
	for _, line := range g.LogMsg {
		logMsg += line + "\n"
	}
	if logMsg != "" {
		tmpEmbed.AddField("게임 로그", logMsg)
	}
	s.ChannelMessageSendEmbed(cid, tmpEmbed.MessageEmbed)
}

// 유저의 진영을 반환하는 함수
func (g *Game) getUserTeam(uId string) string {
	roleName := g.GetRole(uId).String()
	for i := 0; i < len(g.RG); i++ {
		if roleName == g.RG[i].RoleName {
			return g.RG[i].Faction
		}
	}
	return ""
}

// 입력받은 진영에 해당하는 유저가 있는지 확인하는 함수
func (g *Game) userExistsOfThisTeam(team string) bool {
	for _, user := range g.UserList {
		if g.getUserTeam(user.UserID) == team {
			return true
		}
	}
	return false
}

// 각 직업 옆에 표시할 팀마크를 생성하는 함수
// `(backtick) 안에 팀마크를 생성할 경우 제대로 표시되지 않으니 주의 필요
func (g *Game) getTeamMark(role string) string {
	team := ""
	for i := 0; i < len(g.RG); i++ {
		if role == g.RG[i].RoleName {
			team = g.RG[i].Faction
			break
		}
	}
	mark := ""
	switch team {
	case "Villager":
		mark = ":house:"
	case "Werewolf":
		mark = ":wolf:"
	case "Tanner":
		mark = ":coffin:"
	}
	return mark
}

// 게임시작시 버려진 직업인지 확인하는 함수
func (g *Game) isOriDisRole(role Role) bool {
	if FindRoleIdx(role, g.RoleSeq) != -1 && len(g.GetOriRoleUsers(role)) == 0 {
		return true
	}
	return false
}
