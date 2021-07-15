package game

import (
	"github.com/bwmarrin/discordgo"
)

type ActionDoppelganger struct {
	// state 에서 가지고 있는 game
	g         *Game
	cpyRoleID int
	info      *DMInfo
}

// NewActionDoppelganger 는 도플갱어행동 스테이트를 만드는 생성자이다.
func NewActionDoppelganger(g *Game) *ActionDoppelganger {
	ac := &ActionDoppelganger{}
	ac.g = g
	ac.cpyRoleID = -1
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
	switch sdpl.cpyRoleID {
	case -1:
		if g.UserList[num-1].UserID == r.UserID {
			s.ChannelMessageSend(r.ChannelID, "자기 자신을 선택할 수 없습니다.")
			return
		}
		if g.IsProtected(g.UserList[num-1].UserID) {
			s.ChannelMessageSend(r.ChannelID, "`"+g.UserList[num-1].nick+"`은(는) 수호자의 방패로 보호되어 있습니다.")
			return
		}
		s.ChannelMessageDelete(r.ChannelID, sdpl.info.MsgID)
		sdpl.info.Choice <- num
	case (&Seer{}).ID():
		if sdpl.info.Code == 0 {
			if g.UserList[num-1].UserID == r.UserID {
				s.ChannelMessageSend(r.ChannelID, "자기 자신을 선택할 수 없습니다.")
				return
			}
			if g.IsProtected(g.UserList[num-1].UserID) {
				s.ChannelMessageSend(r.ChannelID, "`"+g.UserList[num-1].nick+"`은(는) 수호자의 방패로 보호되어 있습니다.")
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
	case (&Robber{}).ID():
		if sdpl.info.Code == 0 {
			if g.UserList[num-1].UserID == r.UserID {
				s.ChannelMessageSend(r.ChannelID, "자기 자신을 선택할 수 없습니다.")
				return
			}
			if g.IsProtected(g.UserList[num-1].UserID) {
				s.ChannelMessageSend(r.ChannelID, "`"+g.UserList[num-1].nick+"`은(는) 수호자의 방패로 보호되어 있습니다.")
				return
			}
			sdpl.info.Code++
			s.ChannelMessageDelete(r.ChannelID, sdpl.info.MsgID)
			sdpl.info.Choice <- num
		}
	case (&TroubleMaker{}).ID():
		if sdpl.info.Code == 0 {
			if g.UserList[num-1].UserID == r.UserID {
				s.ChannelMessageSend(r.ChannelID, "자기 자신을 선택할 수 없습니다.")
				return
			}
			if g.IsProtected(g.UserList[num-1].UserID) {
				s.ChannelMessageSend(r.ChannelID, "`"+g.UserList[num-1].nick+"`은(는) 수호자의 방패로 보호되어 있습니다.")
				return
			}
			sdpl.info.Code++
			s.ChannelMessageDelete(r.ChannelID, sdpl.info.MsgID)
			sdpl.info.Choice <- num
			sdpl.info.MsgID = role.SendUserSelectGuide(player, g, 1)
		} else if sdpl.info.Code == 1 {
			if g.UserList[num-1].UserID == r.UserID {
				s.ChannelMessageSend(r.ChannelID, "자기 자신을 선택할 수 없습니다.")
				return
			}
			if g.IsProtected(g.UserList[num-1].UserID) {
				s.ChannelMessageSend(r.ChannelID, "`"+g.UserList[num-1].nick+"`은(는) 수호자의 방패로 보호되어 있습니다.")
				return
			}
			sdpl.info.Code++
			s.ChannelMessageDelete(r.ChannelID, sdpl.info.MsgID)
			sdpl.info.Choice <- num
		}
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
	switch sdpl.cpyRoleID {
	case (&Seer{}).ID():
		if curInfo.Code == 0 {
			curInfo.Code++
			s.ChannelMessageDelete(r.ChannelID, curInfo.MsgID)
			curInfo.Choice <- -1
			curInfo.MsgID = role.SendUserSelectGuide(player, g, 1)
		}
	}
}

// PressYesBtn 사용자가 yes 이모티콘을 눌렀을 때 ActionDoppelganger에서 하는 동작
func (sdpl *ActionDoppelganger) PressYesBtn(s *discordgo.Session, r *discordgo.MessageReaction) {
	// do nothing
}

// PressNoBtn 사용자가 No 이모티콘을 눌렀을 때 ActionDoppelganger에서 하는 동작
func (sdpl *ActionDoppelganger) PressNoBtn(s *discordgo.Session, r *discordgo.MessageReaction) {
	// do nothing
}

// PressDirBtn 좌 -1, 우 1 사용자가 좌우 방향 이모티콘을 눌렀을 때 ActionDoppelganger에서 하는 동작
func (sdpl *ActionDoppelganger) PressDirBtn(s *discordgo.Session, r *discordgo.MessageReaction, dir int) {
	// do nothing
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
	role.SendUserSelectGuide(user, g, 0)
	idx := <-sdpl.info.Choice
	tar := &TargetObject{0, user.UserID, g.UserList[idx-1].UserID, -1}
	role.Action(tar, user, g)
	role = g.GetRole(g.UserList[idx-1].UserID)
	sdpl.cpyRoleID = role.ID()
	switch role.String() {
	case (&Seer{}).String():
		seerUserList := g.GetOriRoleUsersWithoutDpl(role)
		for _, user := range seerUserList {
			input := <-sdpl.info.Choice
			if input == -1 {
				tar := &TargetObject{3, "", "", <-sdpl.info.Choice - 1}
				role.Action(tar, user, g)
				role.GenLog(tar, user, g)
			} else {
				tar := &TargetObject{2, g.UserList[input-1].UserID, "", -1}
				role.Action(tar, user, g)
				role.GenLog(tar, user, g)
			}
		}
	case (&Robber{}).String():
		rbUserList := g.GetOriRoleUsersWithoutDpl(role)
		for _, user := range rbUserList {
			input := <-sdpl.info.Choice
			tar := &TargetObject{2, g.UserList[input-1].UserID, "", -1}
			role.Action(tar, user, g)
			role.GenLog(tar, user, g)
		}
	case (&TroubleMaker{}).String():
		tmUserList := g.GetOriRoleUsersWithoutDpl(role)
		for _, user := range tmUserList {
			var input1, input2 int
			for {
				input1 = <-sdpl.info.Choice
				input2 = <-sdpl.info.Choice
				if input1 != input2 {
					break
				}
			}
			tar := &TargetObject{0, g.UserList[input1-1].UserID, g.UserList[input2-1].UserID, -1}
			role.Action(tar, user, g)
			role.GenLog(tar, user, g)
		}
	}
	sdpl.stateFinish()
}

// stateFinish 함수는 ActionDoppelganger state 가 종료될 때 호출되는 메소드이다.
func (sdpl *ActionDoppelganger) stateFinish() {
	sdpl.g.CurState = NewActionInGameGroup(sdpl.g)
	sdpl.g.CurState.InitState()
}

// filterReaction 함수는 각 스테이트에서 보낸 메세지에 리액션 했는지 거르는 함수이다.
// 각 스테이트에서 보낸 메세지의 아이디와 리액션이 온 아이디가 동일한지 확인 및
// 메세지에 리액션 한 것을 지워주어야 한다.
func (sdpl *ActionDoppelganger) filterReaction(s *discordgo.Session, r *discordgo.MessageReaction) bool {
	// do nothing
	return false
}
