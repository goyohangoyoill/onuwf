package game

import (
	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
)

// Game 구조체는 게임 진행을 위한 정보를 담고 있는 스트럭처
type Game struct {
	// 현재 게임이 진행중인 서버의 GID
	GuildID string

	// 현재 게임이 진행중인 채널의 CID
	ChanID string

	// 게임을 생성한 방장의 UID
	MasterID string

	// 현재 게임의 참가자들
	UserList []*User

	// 현재 게임에서 순서대로 추가, 중복제거 된 직업들의 목록
	roleSeq []Role

	// 현재 게임에서 사용중인 사용자에게 보여줄 중복 정렬된 직업들의 목록
	RoleView []Role

	// 현재 게임의 진행시점
	CurState State

	// Role을 User별로 매핑시킨 인덱스 테이블
	// <usage : roleIdxTable[userIdx][roleIdx]>
	roleIdxTable    [][]bool
	oriRoleIdxTable [][]bool

	// 게임에서 버려진 직업 목록
	disRole []Role

	// 게임에서 사용하는 세션
	Session *discordgo.Session

	// 게임 진행 상황을 기록하는 로그 메시지 배열
	LogMsg []string

	// 직업의 대한 소개 및 정보
	RG *RoleGuide

	// 유저 입장시  ID가 전달되는 채널
	UserIDChan chan string
}

// NewGame : Game 스트럭처를 생성하는 생성자,
func NewGame(gid, cid, muid string, rg *RoleGuide, uidChan chan string) (g *Game) {
	g = &Game{}
	g.GuildID = gid
	g.ChanID = cid
	g.MasterID = muid
	g.RG = rg
	g.UserIDChan = uidChan
	g.UserList = make([]*User, 0)
	g.roleSeq = make([]Role, 0)
	g.disRole = make([]Role, 0)
	g.LogMsg = make([]string, 0)
	g.RG = rg
	p := &Prepare{g, 1, nil, nil}
	p.InitEmbed()
	g.CurState = p
	return
}

// SendVoteMsg 는 현재 참가자 모두에게 DM으로 투표 용지를 전송하고,
// 각각의 투표 용지별로 userList index 순서에 맞춰 MsgID 배열을 반환해주는 함수이다.
func (g *Game) SendVoteMsg(s *discordgo.Session) (messageIDs []string) {
	messageIDs = make([]string, len(g.userList))
	for i, me := range g.userList {
		msg := ""
		userListExceptMe := append(g.userList[:i], g.userList[i:]...)
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

// SetUserByID 는 게임에 입장한 유저의 정보를 게임 데이터에 추가하는 함수입니다.
func (g *Game) SetUserByID(uid string) {
	var newOne *User
	newOne.userID = uid
	dgUser, _ := g.Session.User(uid)
	newOne.nick = dgUser.Username
	newOne.chanID = g.ChanID
	uChan, _ := g.Session.UserChannelCreate(uid)
	newOne.dmChanID = uChan.ID
	g.userList = append(g.userList, newOne)
	g.UserIDChan <- uid
}

// DelUserByID 는 입장되어 있는 유저의 정보를 모두 삭제해주는 함수입니다.
func (g *Game) DelUserByID(uid string) {

}

// DelUserByIndex 는 게임에 입장한 유저를 인덱스 번호로 지우는 함수입니다.
func (g *Game) DelUserByIndex(index int) {
	g.userList = append(g.userList[:index], g.userList[index+1:]...)
}

// FindUserByUID UID 로 user 인스턴스를 구하는 함수
func (g *Game) FindUserByUID(uid string) (target *User) {
	for i, item := range g.userList {
		if item.userID == uid {
			return g.userList[i]
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
	loop := len(g.roleSeq)
	idx := FindUserIdx(uid, g.userList)

	for i := 0; i < loop; i++ {
		if g.roleIdxTable[idx][i] {
			return g.roleSeq[i]
		}
	}
	return nil
}

// 유저의 직업을 업데이트
func (g *Game) setRole(uid string, item Role) {
	userIdx := FindUserIdx(uid, g.userList)
	roleIdx := FindRoleIdx(item, g.roleSeq)
	loop := len(g.roleSeq)

	for i := 0; i < loop; i++ {
		g.roleIdxTable[userIdx][i] = false
	}
	g.roleIdxTable[userIdx][roleIdx] = true
}

// SetDisRole 버려진 직업을 업데이트
func (g *Game) SetDisRole(disRoleIdx int, item Role) {
	g.disRole[disRoleIdx] = item
}

// SwapRoleFromUser 두 유저의 직업을 서로 교환
func (g *Game) SwapRoleFromUser(uid1, uid2 string) {
	role1 := g.GetRole(uid1)
	role2 := g.GetRole(uid2)
	g.setRole(uid1, role2)
	g.setRole(uid2, role1)
}

// GetDisRole 버려진 직업 중 하나 확인.
func (g *Game) GetDisRole(disRoleIdx int) Role {
	return g.disRole[disRoleIdx]
}

// SwapRoleFromDiscard 유저 직업과 버려진 직업을 교환.
func (g *Game) SwapRoleFromDiscard(uid string, disRoleIdx int) {
	role1 := g.GetDisRole(disRoleIdx)
	role2 := g.GetRole(uid)
	g.setRole(uid, role1)
	g.SetDisRole(disRoleIdx, role2)
}

// GetRoleUsers 특정 직업의 유저 목록 반환.
func (g *Game) GetRoleUsers(find Role) (users []*User) {
	result := make([]*User, 0)
	loop := len(g.userList)
	idx := FindRoleIdx(find, g.roleSeq)
	for i := 0; i < loop; i++ {
		if g.roleIdxTable[i][idx] {
			result = append(result, g.userList[i])
		}
	}
	return result
}

// RotateAllUserRole  모든 사람들의 직업을 입장순서별로 한칸 회전.
func (g *Game) RotateAllUserRole() {
	loop := len(g.userList)

	tmpRole := g.GetRole(g.userList[loop-1].userID)
	for i := loop - 1; i > 0; i++ {
		item := g.GetRole(g.userList[i-1].userID)
		g.setRole(g.userList[i].userID, item)
	}
	g.setRole(g.userList[0].userID, tmpRole)
}

// SetPower 유저에게 특수권한 부여
func (g *Game) SetPower(power int) {
	// TODO 내부 구현.
}

// CopyRole 특정 유저의 직업을 복사.
func (g *Game) CopyRole(destUID, srcUID string) {
	srcRole := g.GetRole(srcUID)
	g.setRole(destUID, srcRole)
}

// FindUserIdx 유저의 인덱스 찾기를 위한 함수
func FindUserIdx(uid string, target []*User) int {
	for i, item := range target {
		if uid == item.userID {
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
