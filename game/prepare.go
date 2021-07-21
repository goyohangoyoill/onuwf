package game

import (
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
)

type Prepare struct {
	// state에서 가지고 있는 game
	g *Game
	// roleAdd 메세지의 pageNum
	pageNum int
	// 직업추가 확인용 메세지
	RoleAddMsg *discordgo.Message
	// 게임입장 확인용 메세지
	EnterGameMsg *discordgo.Message
}

const pageMax = 3

// PressNumBtn 사용자가 숫자 이모티콘을 눌렀을 때 Prepare에서 하는 동작
func (sPrepare *Prepare) PressNumBtn(s *discordgo.Session, r *discordgo.MessageReaction, num int) {
	// 게임 진행과 관련된 메세지에 달린 리액션 지운다
	if sPrepare.filterReaction(s, r) {
		return
	}
	s.MessageReactionRemove(sPrepare.g.ChanID, r.MessageID, r.Emoji.Name, r.UserID)
	// 직업추가 방장만 가능
	if r.MessageID == sPrepare.RoleAddMsg.ID && r.UserID != sPrepare.g.MasterID {
		s.ChannelMessageSend(sPrepare.g.ChanID, "직업 설정은 방장만 가능합니다")
		return
	}
	num = num + sPrepare.pageNum*pageMax - 1
	if num >= len(sPrepare.g.RG) {
		return
	}
	if num == 2 {
		s.ChannelMessageSend(sPrepare.g.ChanID, "늑대인간은 2개 있어야 합니다")
		return
	}
	// role 생성해서 game의 RoleView와 RoleSeq에 추가
	sPrepare.g.AddRole(num)
	// 메세지 반영
	s.ChannelMessageEditEmbed(sPrepare.g.ChanID, sPrepare.RoleAddMsg.ID, sPrepare.NewRoleAddEmbed().MessageEmbed)
}

// PressDisBtn 사용자가 버려진 카드 이모티콘을 눌렀을 때 Prepare에서 하는 동작
func (sPrepare *Prepare) PressDisBtn(s *discordgo.Session, r *discordgo.MessageReaction) {
}

// PressYesBtn 사용자가 yes 이모티콘을 눌렀을 때 Prepare에서 하는 동작
func (sPrepare *Prepare) PressYesBtn(s *discordgo.Session, r *discordgo.MessageReaction) {
	// 게임 진행과 관련된 메세지에 달린 리액션 지운다
	if sPrepare.filterReaction(s, r) {
		return
	}
	s.MessageReactionRemove(sPrepare.g.ChanID, r.MessageID, r.Emoji.Name, r.UserID)
	//user 생성해서 append()
	sPrepare.g.SetUserByID(r.UserID)
	// 메세지 반영
	s.ChannelMessageEditEmbed(sPrepare.g.ChanID, sPrepare.EnterGameMsg.ID, sPrepare.NewEnterEmbed().MessageEmbed)
	s.ChannelMessageEditEmbed(sPrepare.g.ChanID, sPrepare.RoleAddMsg.ID, sPrepare.NewRoleAddEmbed().MessageEmbed)
}

// PressNoBtn 사용자가 No 이모티콘을 눌렀을 때 Prepare에서 하는 동작
func (sPrepare *Prepare) PressNoBtn(s *discordgo.Session, r *discordgo.MessageReaction) {
	// 게임 진행과 관련된 메세지에 달린 리액션 지운다
	if sPrepare.filterReaction(s, r) {
		return
	}
	s.MessageReactionRemove(sPrepare.g.ChanID, r.MessageID, r.Emoji.Name, r.UserID)
	// userList에서 지우고
	sPrepare.g.DelUserByID(r.UserID)
	// 메세지 반영
	s.ChannelMessageEditEmbed(sPrepare.g.ChanID, sPrepare.EnterGameMsg.ID, sPrepare.NewEnterEmbed().MessageEmbed)
	s.ChannelMessageEditEmbed(sPrepare.g.ChanID, sPrepare.RoleAddMsg.ID, sPrepare.NewRoleAddEmbed().MessageEmbed)
}

// PressDirBtn 좌 -1, 우 1 사용자가 좌우 방향 이모티콘을 눌렀을 때 Prepare에서 하는 동작
func (sPrepare *Prepare) PressDirBtn(s *discordgo.Session, r *discordgo.MessageReaction, dir int) {
	// 게임 진행과 관련된 메세지에 달린 리액션 지운다
	if sPrepare.filterReaction(s, r) {
		return
	}
	s.MessageReactionRemove(sPrepare.g.ChanID, r.MessageID, r.Emoji.Name, r.UserID)
	// 방장만 시작 및 직업추가 가능
	if r.UserID != sPrepare.g.MasterID {
		switch r.MessageID {
		case sPrepare.EnterGameMsg.ID:
			s.ChannelMessageSend(sPrepare.g.ChanID, "게임 시작은 방장만 가능합니다")
		case sPrepare.RoleAddMsg.ID:
			s.ChannelMessageSend(sPrepare.g.ChanID, "직업 설정은 방장만 가능합니다")
		}
		return
	}
	switch r.MessageID {
	// 입장 메세지에서 리액션한거라면
	case sPrepare.EnterGameMsg.ID:
		// 게임 시작
		if dir == 1 && len(sPrepare.g.RoleView) == len(sPrepare.g.UserList)+3 {
			sPrepare.stateFinish()
		}
	// 직업추가 메세지에서 리액션한거라면
	case sPrepare.RoleAddMsg.ID:
		// pageNum 증감
		sPrepare.pageNum += dir
		if sPrepare.pageNum > (len(sPrepare.g.RG) / pageMax) {
			sPrepare.pageNum = 0
		} else if sPrepare.pageNum < 0 {
			sPrepare.pageNum = len(sPrepare.g.RG) / pageMax
		}
		// 직업 추가 메세지 반영
		s.ChannelMessageEditEmbed(sPrepare.g.ChanID, sPrepare.RoleAddMsg.ID, sPrepare.NewRoleAddEmbed().MessageEmbed)
	}
}

// PressBmkBtn DB에 저장된 정보를 load 하는 동작
func (sPrepare *Prepare) PressBmkBtn(s *discordgo.Session, r *discordgo.MessageReaction) {
	if sPrepare.filterReaction(s, r) {
		return
	}
	s.MessageReactionRemove(sPrepare.g.ChanID, r.MessageID, r.Emoji.Name, r.UserID)
	fmt.Println(sPrepare.g.FormerRole)
	fmt.Println(sPrepare.g.MasterID)
	if r.MessageID == sPrepare.RoleAddMsg.ID {
		if r.UserID == sPrepare.g.MasterID && sPrepare.g.FormerRole != nil {
			rLen := len(sPrepare.g.FormerRole)
			for i := 0; i < rLen; i++ {
				sPrepare.g.AddRole(sPrepare.g.FormerRole[i])
			}
		}
	}

	// 입장 확인 메세지 반영
	s.ChannelMessageEditEmbed(sPrepare.g.ChanID, sPrepare.EnterGameMsg.ID, sPrepare.NewEnterEmbed().MessageEmbed)
	// 직업 추가 메세지 반영
	s.ChannelMessageEditEmbed(sPrepare.g.ChanID, sPrepare.RoleAddMsg.ID, sPrepare.NewRoleAddEmbed().MessageEmbed)
}

// InitState 함수는 prepare state가 시작할 때 입장, 직업추가 메세지를 보냅니다.
func (sPrepare *Prepare) InitState() {
	// 늑대인간 2개 추가
	sPrepare.g.AddRole(2)
	sPrepare.g.AddRole(2)
	enterEmbed := sPrepare.NewEnterEmbed()
	roleEmbed := sPrepare.NewRoleAddEmbed()
	s := sPrepare.g.Session
	sPrepare.RoleAddMsg, _ = s.ChannelMessageSendEmbed(sPrepare.g.ChanID, roleEmbed.MessageEmbed)
	nowChan := sPrepare.g.ChanID
	// 직업 추가 메시지에 안내 버튼을 연결
	s.MessageReactionAdd(nowChan, sPrepare.RoleAddMsg.ID, sPrepare.g.Emj["LEFT"])
	s.MessageReactionAdd(nowChan, sPrepare.RoleAddMsg.ID, sPrepare.g.Emj["RIGHT"])
	for i := 0; i < pageMax; i++ {
		s.MessageReactionAdd(sPrepare.RoleAddMsg.ChannelID, sPrepare.RoleAddMsg.ID, sPrepare.g.Emj["n"+strconv.Itoa(i+1)])
	}
	s.MessageReactionAdd(sPrepare.RoleAddMsg.ChannelID, sPrepare.RoleAddMsg.ID, sPrepare.g.Emj["BOOKMARK"])
	sPrepare.EnterGameMsg, _ = s.ChannelMessageSendEmbed(sPrepare.g.ChanID, enterEmbed.MessageEmbed)
	// 게임 입장 메시지에 안내 버튼을 연결
	s.MessageReactionAdd(nowChan, sPrepare.EnterGameMsg.ID, sPrepare.g.Emj["YES"])
	s.MessageReactionAdd(nowChan, sPrepare.EnterGameMsg.ID, sPrepare.g.Emj["NO"])
	s.MessageReactionAdd(nowChan, sPrepare.EnterGameMsg.ID, sPrepare.g.Emj["RIGHT"])
}

func (sPrepare *Prepare) stateFinish() {
	sPrepare.g.CurState = &StartGame{sPrepare.g}
	msg, _ := sPrepare.g.Session.ChannelMessageSend(sPrepare.g.ChanID, "각자의 직업을 배정 중입니다...")
	sPrepare.g.GameStateMID = msg.ID
	sPrepare.g.GameStartedChan <- true
	s := sPrepare.g.Session
	s.ChannelMessageDelete(sPrepare.g.ChanID, sPrepare.EnterGameMsg.ID)
	s.ChannelMessageDelete(sPrepare.g.ChanID, sPrepare.RoleAddMsg.ID)
	sPrepare.g.CurState.InitState()
}

// filterReaction 함수는 입장 메세지랑 직업추가 메세지에 리액션한게 아니면 걸러준다.
func (sPrepare *Prepare) filterReaction(s *discordgo.Session, r *discordgo.MessageReaction) bool {
	// 현재 스테이트에서 보낸 메세지에 리액션한 게 아니면 거름
	if !(r.MessageID == sPrepare.EnterGameMsg.ID || r.MessageID == sPrepare.RoleAddMsg.ID) {
		return true
	}
	if nil == sPrepare.g.FindUserByUID(r.UserID) && !(r.MessageID == sPrepare.EnterGameMsg.ID && r.Emoji.Name == sPrepare.g.Emj["YES"]) {
		return true
	}
	return false
}

// NewRoleAddEmbed 함수는 role guide와 현재 게임에 추가된 직업 / 게임의 참여중인 인원수 + 3 임베드를 만든다
func (sPrepare *Prepare) NewRoleAddEmbed() *embed.Embed {
	roleEmbed := embed.NewEmbed()
	lower := sPrepare.pageNum*pageMax + 1
	upper := (sPrepare.pageNum + 1) * pageMax
	if upper > len(sPrepare.g.RG) {
		upper = len(sPrepare.g.RG)
	}
	roleEmbed.SetTitle("직업 추가" + "（" + strconv.Itoa(sPrepare.pageNum+1) + "/" + strconv.Itoa(len(sPrepare.g.RG)/pageMax+1) + "）")
	msg := ""
	for i := lower - 1; i < upper; i++ {
		if (i)%(pageMax/2) == 0 {
			msg += "\n"
		}
		msg += sPrepare.g.Emj["n"+strconv.Itoa(i%pageMax+1)] + " `" + sPrepare.g.RG[i].RoleName + "`"
		msg += sPrepare.g.getTeamMark(sPrepare.g.RG[i].RoleName)
	}
	roleEmbed.AddField("직업목록", msg)
	roleStr := ""
	if len(sPrepare.g.RoleView) == 0 {
		roleStr += "*추가된 직업이 없습니다.*"
	} else {
		for _, item := range sPrepare.g.RoleSeq {
			cnt := sPrepare.g.RoleCount(item, sPrepare.g.RoleView)
			roleStr += item.String() + " " + strconv.Itoa(cnt) + "개"
			if cnt == sPrepare.g.RG[item.ID()].Max {
				roleStr += " 최대"
			}
			roleStr += "\n"
		}
	}
	roleEmbed.AddField("추가된 직업", roleStr)
	roleEmbed.SetFooter("현재 인원에 맞는 직업 수: " + strconv.Itoa(len(sPrepare.g.RoleView)) + " / " + strconv.Itoa(len(sPrepare.g.UserList)+3))
	return roleEmbed
}

// newEnterEmbed 함수는 게임 참여자 목록 임베드를 만든다
func (sPrepare *Prepare) NewEnterEmbed() *embed.Embed {
	enterEmbed := embed.NewEmbed()
	enterEmbed.SetTitle("게임 참가")
	enterStr := ""
	for _, item := range sPrepare.g.UserList {
		enterStr += "`" + item.nick + "`\n"
	}
	enterEmbed.AddField("참가자 목록", "현재 참가 인원: "+strconv.Itoa(len(sPrepare.g.UserList))+"명\n"+enterStr)
	enterEmbed.SetFooter("(최대 10명, 방장은 나갈 수 없음)\n⭕: 입장 ❌: 퇴장 ▶️: 시작")
	return enterEmbed
}
