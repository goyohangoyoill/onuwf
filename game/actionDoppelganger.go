package game

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

type ActionDoppelganger struct {
	// state 에서 가지고 있는 game
	g             *Game
	cpyRoleString string
	info          *DMInfo
}

// NewActionDoppelganger 는 도플갱어행동 스테이트를 만드는 생성자이다.
func NewActionDoppelganger(g *Game) *ActionDoppelganger {
	ac := &ActionDoppelganger{}
	ac.g = g
	ac.cpyRoleString = (&Doppelganger{}).String()
	ac.info = &DMInfo{"", make(chan int), 0}
	return ac
}

// PressNumBtn 사용자가 숫자 이모티콘을 눌렀을 때 ActionDoppelganger에서 하는 동작
func (sdpl *ActionDoppelganger) PressNumBtn(s *discordgo.Session, r *discordgo.MessageReaction, num int) {
	g := sdpl.g
	role := g.GetOriRole(r.UserID)
	if role.String() != (&Doppelganger{}).String() && !g.IsDoppel(r.UserID) {
		return
	}
	player := g.FindUserByUID(r.UserID)
	switch sdpl.cpyRoleString {
	case (&Doppelganger{}).String():
		if g.UserList[num-1].UserID == r.UserID {
			s.ChannelMessageSend(r.ChannelID, "자기 자신을 선택할 수 없습니다.")
			return
		}
		if g.IsProtected(g.UserList[num-1].UserID) {
			s.ChannelMessageSend(r.ChannelID, "`"+g.UserList[num-1].nick+"`은(는)\n수호자의 방패로 보호되어 있습니다.")
			return
		}
		s.ChannelMessageDelete(r.ChannelID, sdpl.info.MsgID)
		sdpl.info.Choice <- num
	case (&Sentinel{}).String():
		if g.UserList[num-1].UserID == r.UserID {
			s.ChannelMessageSend(r.ChannelID, "자기 자신을 선택할 수 없습니다.")
			return
		}
		sdpl.info.Choice <- num
		s.ChannelMessageDelete(r.ChannelID, sdpl.info.MsgID)
	case (&Seer{}).String():
		if sdpl.info.Code == 0 {
			if g.UserList[num-1].UserID == r.UserID {
				s.ChannelMessageSend(r.ChannelID, "자기 자신을 선택할 수 없습니다.")
				return
			}
			if g.IsProtected(g.UserList[num-1].UserID) {
				s.ChannelMessageSend(r.ChannelID, "`"+g.UserList[num-1].nick+"`은(는)\n수호자의 방패로 보호되어 있습니다.")
				return
			}
			sdpl.info.Code--
			s.ChannelMessageDelete(r.ChannelID, sdpl.info.MsgID)
			sdpl.info.Choice <- num
		}
		if sdpl.info.Code == 1 && num <= 3 {
			sdpl.info.Code++
			s.ChannelMessageDelete(r.ChannelID, sdpl.info.MsgID)
			sdpl.info.Choice <- num
		}
	case (&Robber{}).String():
		if sdpl.info.Code == 0 {
			if g.UserList[num-1].UserID == r.UserID {
				s.ChannelMessageSend(r.ChannelID, "자기 자신을 선택할 수 없습니다.")
				return
			}
			if g.IsProtected(g.UserList[num-1].UserID) {
				s.ChannelMessageSend(r.ChannelID, "`"+g.UserList[num-1].nick+"`은(는)\n수호자의 방패로 보호되어 있습니다.")
				return
			}
			sdpl.info.Code++
			s.ChannelMessageDelete(r.ChannelID, sdpl.info.MsgID)
			sdpl.info.Choice <- num
		}
	case (&TroubleMaker{}).String():
		if sdpl.info.Code == 0 {
			if g.UserList[num-1].UserID == r.UserID {
				s.ChannelMessageSend(r.ChannelID, "자기 자신을 선택할 수 없습니다.")
				return
			}
			if g.IsProtected(g.UserList[num-1].UserID) {
				s.ChannelMessageSend(r.ChannelID, "`"+g.UserList[num-1].nick+"`은(는)\n수호자의 방패로 보호되어 있습니다.")
				return
			}
			sdpl.info.Code++
			s.ChannelMessageDelete(r.ChannelID, sdpl.info.MsgID)
			sdpl.info.MsgID = role.SendUserSelectGuide(player, g, 1)
			sdpl.info.Choice <- num
		} else if sdpl.info.Code == 1 {
			if g.UserList[num-1].UserID == r.UserID {
				s.ChannelMessageSend(r.ChannelID, "자기 자신을 선택할 수 없습니다.")
				return
			}
			if g.IsProtected(g.UserList[num-1].UserID) {
				s.ChannelMessageSend(r.ChannelID, "`"+g.UserList[num-1].nick+"`은(는)\n수호자의 방패로 보호되어 있습니다.")
				return
			}
			sdpl.info.Code++
			s.ChannelMessageDelete(r.ChannelID, sdpl.info.MsgID)
			sdpl.info.Choice <- num
		}
	case (&Drunk{}).String():
		s.ChannelMessageDelete(r.ChannelID, sdpl.info.MsgID)
		sdpl.info.Choice <- num
	}
}

// PressDisBtn 사용자가 버려진 카드 이모티콘을 눌렀을 때 ActionDoppelganger에서 하는 동작
func (sdpl *ActionDoppelganger) PressDisBtn(s *discordgo.Session, r *discordgo.MessageReaction) {
	g := sdpl.g
	role := g.GetOriRole(r.UserID)
	player := g.FindUserByUID(r.UserID)
	if !(g.IsDoppel(r.UserID)) {
		return
	}
	curInfo := sdpl.info
	switch sdpl.cpyRoleString {
	case (&Seer{}).String():
		if curInfo.Code == 0 {
			curInfo.Code++
			s.ChannelMessageDelete(r.ChannelID, curInfo.MsgID)
			curInfo.Choice <- -1
			curInfo.MsgID = role.SendUserSelectGuide(player, g, 1)
		}
	case (&Sentinel{}).String():
		if sdpl.info.MsgID != r.MessageID {
			return
		}
		s.ChannelMessageDelete(r.ChannelID, sdpl.info.MsgID)
		s.ChannelMessageSend(r.ChannelID, "아무도 방패로 보호하지 않았습니다.")
		sdpl.info.Choice <- -1
	}
}

// PressYesBtn 사용자가 yes 이모티콘을 눌렀을 때 ActionDoppelganger에서 하는 동작
func (sdpl *ActionDoppelganger) PressYesBtn(s *discordgo.Session, r *discordgo.MessageReaction) {
}

// PressNoBtn 사용자가 No 이모티콘을 눌렀을 때 ActionDoppelganger에서 하는 동작
func (sdpl *ActionDoppelganger) PressNoBtn(s *discordgo.Session, r *discordgo.MessageReaction) {
}

// PressDirBtn 좌 -1, 우 1 사용자가 좌우 방향 이모티콘을 눌렀을 때 ActionDoppelganger에서 하는 동작
func (sdpl *ActionDoppelganger) PressDirBtn(s *discordgo.Session, r *discordgo.MessageReaction, dir int) {
	// do nothing
}

// PressBmkBtn DB에 저장된 정보를 load 하는 동작
func (sdpl *ActionDoppelganger) PressBmkBtn(s *discordgo.Session, r *discordgo.MessageReaction) {
	//do nothing
}

// InitState 함수는 ActionDoppelganger state 가 시작되었을 때 호출되는 메소드이다.
func (sdpl *ActionDoppelganger) InitState() {
	g := sdpl.g
	var role Role
	role = &Doppelganger{}
	if FindRoleIdx(role, g.RoleSeq) == -1 {
		sdpl.stateFinish()
		return
	}
	dplUserList := g.GetOriRoleUsers(role)
	if len(dplUserList) == 0 {
		sdpl.stateFinish()
		return
	}
	user := dplUserList[0]
	sdpl.info.MsgID = role.SendUserSelectGuide(user, g, 0)
	idx := <-sdpl.info.Choice - 1
	tar := &TargetObject{0, user.UserID, g.UserList[idx].UserID, -1}
	role.GenLog(tar, user, g)
	role.Action(tar, user, g)
	role = g.GetRole(g.UserList[idx].UserID)
	if role.String() == (&Werewolf{}).String() {
		sdpl.stateFinish()
		return
	}
	sdpl.info.MsgID = role.SendUserSelectGuide(user, g, 0)
	sdpl.info.Code = 0
	sdpl.cpyRoleString = role.String()
	switch role.String() {
	case (&Seer{}).String():
		input := <-sdpl.info.Choice
		if input == -1 {
			tar := &TargetObject{3, "", "", <-sdpl.info.Choice - 1}
			role.GenLog(tar, user, g)
			role.Action(tar, user, g)
		} else {
			tar := &TargetObject{2, g.UserList[input-1].UserID, "", -1}
			role.GenLog(tar, user, g)
			role.Action(tar, user, g)
		}
	case (&Robber{}).String():
		input := <-sdpl.info.Choice
		tar := &TargetObject{2, g.UserList[input-1].UserID, "", -1}
		role.GenLog(tar, user, g)
		role.Action(tar, user, g)
	case (&TroubleMaker{}).String():
		var input1, input2 int
		for {
			input1 = <-sdpl.info.Choice
			input2 = <-sdpl.info.Choice
			if input1 != input2 {
				break
			}
			g.Session.ChannelMessageSend(user.dmChanID, "같은 사람을 두 번 선택할 수 없습니다")
			sdpl.info.MsgID = role.SendUserSelectGuide(user, g, 0)
			sdpl.info.Code = 0
		}
		tar := &TargetObject{0, g.UserList[input1-1].UserID, g.UserList[input2-1].UserID, -1}
		role.GenLog(tar, user, g)
		role.Action(tar, user, g)
	case (&Drunk{}).String():
		input := <-sdpl.info.Choice
		tar := &TargetObject{1, user.UserID, "", input - 1}
		role.Action(tar, user, g)
		role.GenLog(tar, user, g)
	case (&Sentinel{}).String():
		input := <-sdpl.info.Choice
		if input == -1 {
			tar := &TargetObject{2, "", "", -1}
			role.GenLog(tar, user, g)
		} else {
			tar := &TargetObject{2, g.UserList[input-1].UserID, "", -1}
			role.Action(tar, user, g)
			role.GenLog(tar, user, g)
		}
	}
	sdpl.stateFinish()
}

// stateFinish 함수는 ActionDoppelganger state 가 종료될 때 호출되는 메소드이다.
func (sdpl *ActionDoppelganger) stateFinish() {
	// 도플갱어카드가 처음에 버려졌을 경우 다음 스테이트 넘어가기전 딜레이를 걸어줌
	// DM 날라오는 타이밍을 통해 도플갱어인 유저가 있는지 유추하지 못하도록 하기 위함
	if sdpl.g.isOriDisRole(&Doppelganger{}) {
		delaySec := sdpl.g.config.DisRoleDelaySec
		time.Sleep(time.Duration(delaySec) * time.Second)
	}
	sdpl.g.CurState = NewActionInGameGroup(sdpl.g)
	sdpl.g.CurState.InitState()
}

// filterReaction 함수는 각 스테이트에서 보낸 메세지에 리액션 했는지 거르는 함수이다.
// 각 스테이트에서 보낸 메세지의 아이디와 리액션이 온 아이디가 동일한지 확인 및
// 메세지에 리액션 한 것을 지워주어야 한다.
func (sdpl *ActionDoppelganger) filterReaction(s *discordgo.Session, r *discordgo.MessageReaction) bool {
	return false
}
