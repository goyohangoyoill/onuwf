package game

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

type ActionInGameGroup struct {
	// state 에서 가지고 있는 game
	g *Game

	Info map[string]*DMInfo
}

type DMInfo struct {
	MsgID  string
	Choice chan int
	Code   int
}

func NewActionInGameGroup(g *Game) *ActionInGameGroup {
	ac := &ActionInGameGroup{}
	ac.g = g
	ac.Info = make(map[string]*DMInfo)
	return ac
}

// PressNumBtn 사용자가 숫자 이모티콘을 눌렀을 때 ActionInGameGroup에서 하는 동작
func (sActionInGameGroup *ActionInGameGroup) PressNumBtn(s *discordgo.Session, r *discordgo.MessageReaction, num int) {
	if FindUserIdx(r.UserID, sActionInGameGroup.g.UserList) == -1 {
		return
	}
	role := sActionInGameGroup.g.GetOriRole(r.UserID)
	player := sActionInGameGroup.g.FindUserByUID(r.UserID)
	curInfo := sActionInGameGroup.Info[player.UserID]
	switch role.String() {
	case (&Werewolf{}).String():
		if curInfo.Code == 1 && num <= 3 {
			curInfo.Code++
			s.ChannelMessageDelete(r.ChannelID, curInfo.MsgID)
			curInfo.Choice <- num
		}
	case (&Seer{}).String():
		if curInfo.Code == 0 {
			if sActionInGameGroup.g.UserList[num-1].UserID == r.UserID {
				s.ChannelMessageSend(r.ChannelID, "자기 자신을 선택할 수 없습니다.")
				return
			}
			if sActionInGameGroup.g.IsProtected(sActionInGameGroup.g.UserList[num-1].UserID) {
				s.ChannelMessageSend(r.ChannelID, "`"+sActionInGameGroup.g.UserList[num-1].nick+"`은(는) 수호자의 방패로 보호되어 있습니다.")
				return
			}
			curInfo.Code--
			s.ChannelMessageDelete(r.ChannelID, curInfo.MsgID)
			curInfo.Choice <- num
		}
		if curInfo.Code == 1 && num <= 3 {
			curInfo.Code++
			s.ChannelMessageDelete(r.ChannelID, curInfo.MsgID)
			curInfo.Choice <- num
		}
	case (&Robber{}).String():
		if curInfo.Code == 0 {
			if sActionInGameGroup.g.UserList[num-1].UserID == r.UserID {
				s.ChannelMessageSend(r.ChannelID, "자기 자신을 선택할 수 없습니다.")
				return
			}
			if sActionInGameGroup.g.IsProtected(sActionInGameGroup.g.UserList[num-1].UserID) {
				s.ChannelMessageSend(r.ChannelID, "`"+sActionInGameGroup.g.UserList[num-1].nick+"`은(는) 수호자의 방패로 보호되어 있습니다.")
				return
			}
			curInfo.Code++
			s.ChannelMessageDelete(r.ChannelID, curInfo.MsgID)
			curInfo.Choice <- num
		}
	case (&TroubleMaker{}).String():
		if curInfo.Code == 0 {
			if sActionInGameGroup.g.UserList[num-1].UserID == r.UserID {
				s.ChannelMessageSend(r.ChannelID, "자기 자신을 선택할 수 없습니다.")
				return
			}
			if sActionInGameGroup.g.IsProtected(sActionInGameGroup.g.UserList[num-1].UserID) {
				s.ChannelMessageSend(r.ChannelID, "`"+sActionInGameGroup.g.UserList[num-1].nick+"`은(는) 수호자의 방패로 보호되어 있습니다.")
				return
			}
			curInfo.Code++
			s.ChannelMessageDelete(r.ChannelID, curInfo.MsgID)
			curInfo.Choice <- num
			curInfo.MsgID = role.SendUserSelectGuide(player, sActionInGameGroup.g, 1)
		} else if curInfo.Code == 1 {
			if sActionInGameGroup.g.UserList[num-1].UserID == r.UserID {
				s.ChannelMessageSend(r.ChannelID, "자기 자신을 선택할 수 없습니다.")
				return
			}
			if sActionInGameGroup.g.IsProtected(sActionInGameGroup.g.UserList[num-1].UserID) {
				s.ChannelMessageSend(r.ChannelID, "`"+sActionInGameGroup.g.UserList[num-1].nick+"`은(는) 수호자의 방패로 보호되어 있습니다.")
				return
			}
			curInfo.Code++
			s.ChannelMessageDelete(r.ChannelID, curInfo.MsgID)
			curInfo.Choice <- num
		}
	}
}

// PressDisBtn 사용자가 버려진 카드 이모티콘을 눌렀을 때 ActionInGameGroup에서 하는 동작
func (sActionInGameGroup *ActionInGameGroup) PressDisBtn(s *discordgo.Session, r *discordgo.MessageReaction) {
	if FindUserIdx(r.UserID, sActionInGameGroup.g.UserList) == -1 {
		return
	}
	role := sActionInGameGroup.g.GetOriRole(r.UserID)
	player := sActionInGameGroup.g.FindUserByUID(r.UserID)
	curInfo := sActionInGameGroup.Info[player.UserID]
	switch role.String() {
	case (&Seer{}).String():
		if curInfo.Code == 0 {
			curInfo.Code++
			s.ChannelMessageDelete(r.ChannelID, curInfo.MsgID)
			curInfo.MsgID = role.SendUserSelectGuide(player, sActionInGameGroup.g, 1)
			curInfo.Choice <- -1
		}
	}
}

// PressYesBtn 사용자가 yes 이모티콘을 눌렀을 때 ActionInGameGroup에서 하는 동작
func (sActionInGameGroup *ActionInGameGroup) PressYesBtn(s *discordgo.Session, r *discordgo.MessageReaction) {
}

// PressNoBtn 사용자가 No 이모티콘을 눌렀을 때 ActionInGameGroup에서 하는 동작
func (sActionInGameGroup *ActionInGameGroup) PressNoBtn(s *discordgo.Session, r *discordgo.MessageReaction) {
}

// PressDirBtn 좌 -1, 우 1 사용자가 좌우 방향 이모티콘을 눌렀을 때 ActionInGameGroup에서 하는 동작
func (sActionInGameGroup *ActionInGameGroup) PressDirBtn(s *discordgo.Session, r *discordgo.MessageReaction, dir int) {
}

// PressBmkBtn DB에 저장된 정보를 load 하는 동작
func (sActionInGameGroup *ActionInGameGroup) PressBmkBtn(s *discordgo.Session, r *discordgo.MessageReaction) {
	//do nothing
}

// InitState 함수는 ActionInGameGroup state 가 시작되었을 때 호출되는 메소드이다.
func (sActionInGameGroup *ActionInGameGroup) InitState() {
	g := sActionInGameGroup.g
	// 늑대인간부터 말썽쟁이까지
	for i := 2; i < 8; i++ {
		curInfo := &DMInfo{"", make(chan int), 0}
		role := GenerateRole(i)
		if role == nil {
			continue
		}
		// 밑에서 getRoleUsers에서 nil나와서 검사 해야됨
		rIdx := FindRoleIdx(role, g.RoleSeq)
		if rIdx == -1 {
			continue
		}
		for _, user := range g.GetOriRoleUsersWithoutDpl(role) {
			sActionInGameGroup.Info[user.UserID] = curInfo
			if role.String() == (&Werewolf{}).String() {
				wolves := g.GetOriRoleUsers(&Werewolf{})
				//wolves = append(wolves, g.GetOriRoleUsers(&Misticwolf{})...)
				//wolves = append(wolves, g.GetOriRoleUsers(&Alphawolf{})...)
				//wolves = append(wolves, g.GetOriRoleUsers(&Dreamwolf{})...)
				if len(wolves) == 1 {
					(sActionInGameGroup.Info[user.UserID]).Code = 1
					curInfo.MsgID = role.SendUserSelectGuide(user, g, 0)
				}
				continue
			}
			curInfo.MsgID = role.SendUserSelectGuide(user, g, 0)
		}
	}
	curInfo := sActionInGameGroup.Info
	for i := 0; i < len(g.RoleSeq); i++ {
		role := g.RoleSeq[i]
		uList := g.GetOriRoleUsers(role)
		if len(uList) == 0 {
			continue
		}
		switch role.String() {
		case (&Werewolf{}).String():
			wfUserList := uList
			onlyWF := wfUserList[0]
			if curInfo[onlyWF.UserID].Code == 1 {
				input := <-curInfo[onlyWF.UserID].Choice
				tar := &TargetObject{3, "", "", input - 1}
				role.Action(tar, onlyWF, g)
				role.GenLog(tar, onlyWF, g)
			} else {
				tar := &TargetObject{-1, "", "", -1}
				role.GenLog(tar, onlyWF, g)
				for _, user := range wfUserList {
					role.Action(tar, user, g)
				}
			}
		case (&Minion{}).String():
			minUserList := uList
			for _, user := range minUserList {
				tar := &TargetObject{-1, "", "", -1}
				role.Action(tar, user, g)
				role.GenLog(tar, user, g)
			}
		case (&Freemason{}).String():
			frmUserList := uList
			tar := &TargetObject{-1, "", "", -1}
			if len(frmUserList) != 0 {
				role.GenLog(tar, frmUserList[0], g)
			}
			for _, user := range frmUserList {
				role.Action(tar, user, g)
			}
		case (&Seer{}).String():
			seerUserList := g.GetOriRoleUsersWithoutDpl(role)
			for _, user := range seerUserList {
				input := <-curInfo[user.UserID].Choice
				if input == -1 {
					tar := &TargetObject{3, "", "", <-curInfo[user.UserID].Choice - 1}
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
				input := <-curInfo[user.UserID].Choice
				tar := &TargetObject{2, g.UserList[input-1].UserID, "", -1}
				role.GenLog(tar, user, g)
				role.Action(tar, user, g)
			}
		case (&TroubleMaker{}).String():
			tmUserList := g.GetOriRoleUsersWithoutDpl(role)
			for _, user := range tmUserList {
				var input1, input2 int
				for {
					input1 = <-curInfo[user.UserID].Choice
					input2 = <-curInfo[user.UserID].Choice
					if input1 != input2 {
						break
					}
					g.Session.ChannelMessageSend(user.dmChanID, "같은 사람을 두 번 선택할 수 없습니다.")
					curInfo[user.UserID].MsgID = role.SendUserSelectGuide(user, g, 0)
					curInfo[user.UserID].Code = 0
				}
				tar := &TargetObject{0, g.UserList[input1-1].UserID, g.UserList[input2-1].UserID, -1}
				role.Action(tar, user, g)
				role.GenLog(tar, user, g)
			}
		}
	}
	sActionInGameGroup.stateFinish()
}

// stateFinish 함수는 ActionInGameGroup state 가 종료될 때 호출되는 메소드이다.
func (sActionInGameGroup *ActionInGameGroup) stateFinish() {
	time.Sleep(time.Second * 3)
	sActionInGameGroup.g.CurState = NewStateBeforeVote(sActionInGameGroup.g)
	sActionInGameGroup.g.CurState.InitState()
}

// filterReaction 함수는 각 스테이트에서 보낸 메세지에 리액션 했는지 거르는 함수이다.
// 각 스테이트에서 보낸 메세지의 아이디와 리액션이 온 아이디가 동일한지 확인 및
// 메세지에 리액션 한 것을 지워주어야 한다.
func (sActionInGameGroup *ActionInGameGroup) filterReaction(s *discordgo.Session, r *discordgo.MessageReaction) bool {
	return false
}
