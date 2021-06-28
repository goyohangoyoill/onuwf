// +build linux,amd64,go1.15,!cgo

package internal

import (
	"os"

	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
)

// Game 구조체는 게임 진행을 위한 정보를 담고 있는 스트럭처
type Game struct {
	// 현재 게임이 진행중인 서버의 GID
	guildID string

	// 현재 게임이 진행중인 채널의 CID
	chanID string

	enterGameMsgID string
	roleAddMsgID   string
	// 게임을 생성한 방장의 UID
	masterID string

	// 현재 게임의 참가자들
	userList []user

	// 현재 게임에서 순서대로 추가, 중복제거 된 직업들의 목록
	roleSeq []role

	// 현재 게임에서 사용중인 사용자에게 보여줄 중복 정렬된 직업들의 목록
	roleView []role

	// 현재 게임의 진행시점
	curState state

	// Role을 User별로 매핑시킨 인덱스 테이블
	// <usage : roleIdxTable[userIdx][roleIdx]>
	roleIdxTable    [][]bool
	oriRoleIdxTable [][]bool

	// 게임에서 버려진 직업 목록
	disRole []role

	// 게임 진행 상황을 기록하는 로그 메시지 배열
	logMsg []string

	killChan chan os.Signal
}

func NewGame(gid, cid, muid string) (g *Game) {
	g = &Game{}
	g.guildID = gid
	g.chanID = cid
	g.masterID = muid
	g.userList = make([]user, 0)
	g.roleSeq = make([]role, 0)
	g.disRole = make([]role, 0)
	g.curState = &StatePrepare{g, 1, nil, nil}
	g.logMsg = make([]string, 0)
	return
}

// SendVoteMsg() 는 현재 참가자 모두에게 DM으로 투표 용지를 전송하고,
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

func (g *Game) SetUserByID(s *discordgo.Session, uid string) {
	var newOne user
	newOne.userID = uid
	dgUser, _ := s.User(uid)
	newOne.nick = dgUser.Username
	newOne.chanID = g.chanID
	uChan, _ := s.UserChannelCreate(uid)
	newOne.dmChanID = uChan.ID
	g.userList = append(g.userList, newOne)
}

// UID 로 user 인스턴스를 구하는 함수
func (g *Game) FindUserByUID(uid string) (target *user) {
	for i, item := range g.userList {
		if item.userID == uid {
			return &g.userList[i]
		}
	}
	return nil
}

// 게임 로그에 메시지를 쌓는 함수.
func (g *Game) AppendLog(msg string) {
	if g.logMsg == nil {
		g.logMsg = make([]string, 0)
	}
	g.logMsg = append(g.logMsg, msg)
}

// 유저의 직업을 반환
func (g *Game) GetRole(uid string) role {
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
func (g *Game) setRole(uid string, item role) {
	userIdx := FindUserIdx(uid, g.userList)
	roleIdx := FindRoleIdx(item, g.roleSeq)
	loop := len(g.roleSeq)

	for i := 0; i < loop; i++ {
		g.roleIdxTable[userIdx][i] = false
	}
	g.roleIdxTable[userIdx][roleIdx] = true
}

// 버려진 직업을 업데이트
func (g *Game) setDisRole(disRoleIdx int, item role) {
	g.disRole[disRoleIdx] = item
}

// 두 유저의 직업을 서로 교환
func (g *Game) SwapRoleFromUser(uid1, uid2 string) {
	role1 := g.getRole(uid1)
	role2 := g.getRole(uid2)
	g.setRole(uid1, role2)
	g.setRole(uid2, role1)
}

// 버려진 직업 중 하나 확인.
func (g *Game) GetDisRole(disRoleIdx int) role {
	return g.disRole[disRoleIdx]
}

// 유저 직업과 버려진 직업을 교환.
func (g *Game) SwapRoleFromDiscard(uid string, disRoleIdx int) {
	role1 := g.getDisRole(disRoleIdx)
	role2 := g.getRole(uid)
	g.setRole(uid, role1)
	g.setDisRole(disRoleIdx, role2)
}

// 특정 직업의 유저 목록 반환.
func (g *Game) GetRoleUsers(find role) (users []user) {
	result := make([]user, 0)
	loop := len(g.userList)

	idx := FindRoleIdx(find, g.roleSeq)

	for i := 0; i < loop; i++ {
		if g.roleIdxTable[i][idx] {
			result = append(result, g.userList[i])
		}
	}

	return result
}

// 모든 사람들의 직업을 입장순서별로 한칸 회전.
func (g *Game) RotateAllUserRole() {
	loop := len(g.userList)

	tmpRole := g.getRole(g.userList[loop-1].userID)
	for i := loop - 1; i > 0; i++ {
		item := g.getRole(g.userList[i-1].userID)
		g.setRole(g.userList[i].userID, item)
	}
	g.setRole(g.userList[0].userID, tmpRole)
}

// 유저에게 특수권한 부여
func (g *Game) SetPower(power int) {
	// TODO 내부 구현.
}

// 특정 유저의 직업을 복사.
func (g *Game) CopyRole(destUID, srcUID string) {
	srcRole := g.getRole(srcUID)
	g.setRole(destUID, srcRole)
}

// 유저의 인덱스 찾기를 위한 함수
func FindUserIdx(uid string, target []user) int {
	for i, item := range target {
		if uid == item.userID {
			return i
		}
	}
	return -1
}

// 직업의 인덱스 찾기를 위한 함수
func FindRoleIdx(r role, target []role) int {
	for i, item := range target {
		if r.String() == item.String() {
			return i
		}
	}
	return -1
}
