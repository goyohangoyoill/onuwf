package game

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
)

type Prepare struct {
	// state에서 가지고 있는 game
	g *Game

	// factory 에서 쓰이게 될 role index
	roleIndex int

	// 직업추가 확인용 메세지
	roleAddMsg *discordgo.Message

	// 게임입장 확인용 메세지
	enterGameMsg *discordgo.Message
}

// PressNumBtn 사용자가 숫자 이모티콘을 눌렀을 때 Prepare에서 하는 동작
func (sPrepare *Prepare) PressNumBtn(s *discordgo.Session, r *discordgo.MessageReactionAdd, num int) {
	// do nothing
}

// PressDisBtn 사용자가 버려진 카드 이모티콘을 눌렀을 때 Prepare에서 하는 동작
func (sPrepare *Prepare) PressDisBtn(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// do nothing
}

// PressYesBtn 사용자가 yes 이모티콘을 눌렀을 때 Prepare에서 하는 동작
func (sPrepare *Prepare) PressYesBtn(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// 입장 메세지에서 리액션한거라면
	if r.MessageID == sPrepare.enterGameMsg.ID {
		//user 생성해서 append()
		sPrepare.g.SetUserByID(r.UserID)
		// 입장 확인 메세지 반영
		s.ChannelMessageEditEmbed(sPrepare.g.ChanID, sPrepare.enterGameMsg.ID, sPrepare.NewEnterEmbed().MessageEmbed)
		// 직업추가 메세지에서 리액션한거라면
	} else if r.MessageID == sPrepare.roleAddMsg.ID {
		// roleFactory에서 현재 roleindex 위치 값을 받아 role 생성
		roleToAdd := GenerateRole(sPrepare.roleIndex)
		// 추가된 role 개수가 max보다 작을 때만 추가
		if len(sPrepare.g.GetRoleUsers(roleToAdd)) < sPrepare.g.RG[sPrepare.roleIndex].Max {
			// RoleView는 ununique sorted니까 RoleView에 중복된 상태로 sort index 찾아서 삽입
			for i, item := range sPrepare.g.RoleView {
				if item.String() <= roleToAdd.String() {
					tmp := append(sPrepare.g.RoleView[:i], roleToAdd)
					sPrepare.g.RoleView = append(tmp, sPrepare.g.RoleView[i:]...)
					break
				}
			}
			// RoleSeq는 unique unsorted니까 RoleSeq에 없으면 append
			if FindRoleIdx(roleToAdd, sPrepare.g.RoleSeq) == -1 {
				sPrepare.g.RoleSeq = append(sPrepare.g.RoleSeq, roleToAdd)
			}
			// 직업 추가 메세지 반영
			s.ChannelMessageEditEmbed(sPrepare.g.ChanID, sPrepare.roleAddMsg.ID, sPrepare.NewRoleEmbed().MessageEmbed)
		}
	}
}

// PressNoBtn 사용자가 No 이모티콘을 눌렀을 때 Prepare에서 하는 동작
func (sPrepare *Prepare) PressNoBtn(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// 입장 메세지에서 리액션한거라면
	if r.MessageID == sPrepare.enterGameMsg.ID {
		// userList에서 지우고
		sPrepare.g.DelUserByID(r.UserID)
		// 입장 확인 메세지에서 지우기
		msg := sPrepare.enterGameMsg.Embeds[0].Description
		strings.Replace(msg, r.UserID+" 입장\n", "", 1)
		sPrepare.enterGameMsg.Embeds[0].Description = msg
		s.ChannelMessageEditEmbed(sPrepare.g.ChanID, sPrepare.enterGameMsg.ID, sPrepare.enterGameMsg.Embeds[0])
		// 직업추가 메세지에서 리액션한거라면
	} else if r.MessageID == sPrepare.roleAddMsg.ID {
		// roleFactory에서 현재 roleindex 위치 값을 받아 role 생성
		roleToRemove := GenerateRole(sPrepare.roleIndex)
		// RoleView는 ununique sorted니까 첫번째로 나오는거 찾아서 지우기
		if i := FindRoleIdx(roleToRemove, sPrepare.g.RoleView); i != -1 {
			sPrepare.g.RoleView = append(sPrepare.g.RoleView[:i], sPrepare.g.RoleView[i+1:]...)
			// RoleSeq는 unique unsorted니까 방금 지운 RoleView에 없으면 지우기
		} else if i := FindRoleIdx(roleToRemove, sPrepare.g.RoleSeq); i != -1 {
			sPrepare.g.RoleSeq = append(sPrepare.g.RoleSeq[:i], sPrepare.g.RoleSeq[i+1:]...)
		}
		// 직업 확인 메세지 보낸 적 있으면 수정하거나 지우기
		if sPrepare.roleAddMsg != nil {
			var msg string
			msg = sPrepare.roleAddMsg.Embeds[0].Description

			rgxstr := regexp.MustCompile(roleToRemove.String() + " " + `\d+.*\n`)
			rgxnum := regexp.MustCompile(`\d`)
			// 직업 확인에 직업이 있을 때
			if rgxstr.MatchString(msg) {
				// 정규표현식으로 "직업 숫자\n" line으로 찾음
				line := rgxstr.FindString(msg)
				// line에서 숫자만 뽑아냄
				num, _ := strconv.Atoi(rgxnum.FindString(line))
				// role guide에 있는 직업의 max랑 같으면
				if num == sPrepare.g.RG[sPrepare.roleIndex].Max {
					// 최대 지우기
					msg = strings.Replace(msg, line, strings.Replace(line, "최대", "", 1), 1)
					// 1이면
				} else if num == 1 {
					msg = strings.Replace(msg, line, "", 1)
				} else {
					// 숫자 감소
					msg = strings.Replace(msg, line, roleToRemove.String()+" "+fmt.Sprint(num-1)+"\n", 1)
				}
			}

			sPrepare.roleAddMsg.Embeds[0].Description = msg
			s.ChannelMessageEditEmbed(sPrepare.g.ChanID, sPrepare.roleAddMsg.ID, sPrepare.roleAddMsg.Embeds[0])
		}
	}
}

// PressDirBtn 좌 -1, 우 1 사용자가 좌우 방향 이모티콘을 눌렀을 때 Prepare에서 하는 동작
func (sPrepare *Prepare) PressDirBtn(s *discordgo.Session, r *discordgo.MessageReactionAdd, dir int) {
	if r.MessageID == sPrepare.enterGameMsg.ID {
		// 게임 시작
		if dir == 1 {
			if len(sPrepare.g.RoleSeq) == len(sPrepare.g.UserList)+3 {
				sPrepare.g.CurState = &Playable{sPrepare.g, sPrepare.roleIndex, sPrepare.roleAddMsg, sPrepare.enterGameMsg}
				s.ChannelMessageSendEmbed(sPrepare.g.ChanID, embed.NewGenericEmbed("게임시작", ""))
			} else {
				s.ChannelMessageSendEmbed(sPrepare.g.ChanID, embed.NewGenericEmbed("인원수 안 맞음", ""))
			}
		}
	} else if r.MessageID == sPrepare.roleAddMsg.ID {
		// roleindex 증감
		sPrepare.roleIndex += dir
		if dir > len(sPrepare.g.RG) {
			dir = 0
		} else if dir <= 0 {
			dir = len(sPrepare.g.RG)
		}
		s.ChannelMessageEditEmbed(sPrepare.g.ChanID, sPrepare.roleAddMsg.ID,
			embed.NewGenericEmbed(sPrepare.g.RG[sPrepare.roleIndex].RoleName, strings.Join(sPrepare.g.RG[sPrepare.roleIndex].RoleGuide, "\n")))
	}
}

// InitEmbed 함수는 게임이 시작할 때 입장, 직업추가 메세지를 보냅니다.
func (sPrepare *Prepare) InitEmbed() {
	enterEmbed := sPrepare.NewEnterEmbed()
	roleEmbed := sPrepare.NewRoleEmbed()
	s := sPrepare.g.Session
	sPrepare.enterGameMsg, _ = s.ChannelMessageSendEmbed(sPrepare.g.ChanID, enterEmbed.MessageEmbed)
	// 게임 입장 메시지에 안내 버튼을 연결
	s.MessageReactionAdd(sPrepare.enterGameMsg.ChannelID, sPrepare.enterGameMsg.ID, sPrepare.g.Emj["YES"])
	s.MessageReactionAdd(sPrepare.enterGameMsg.ChannelID, sPrepare.enterGameMsg.ID, sPrepare.g.Emj["NO"])
	sPrepare.roleAddMsg, _ = s.ChannelMessageSendEmbed(sPrepare.g.ChanID, roleEmbed.MessageEmbed)
	// 직업 추가 메시지에 안내 버튼을 연결
	s.MessageReactionAdd(sPrepare.roleAddMsg.ChannelID, sPrepare.roleAddMsg.ID, sPrepare.g.Emj["YES"])
	s.MessageReactionAdd(sPrepare.roleAddMsg.ChannelID, sPrepare.roleAddMsg.ID, sPrepare.g.Emj["NO"])
	s.MessageReactionAdd(sPrepare.roleAddMsg.ChannelID, sPrepare.roleAddMsg.ID, sPrepare.g.Emj["LEFT"])
	s.MessageReactionAdd(sPrepare.roleAddMsg.ChannelID, sPrepare.roleAddMsg.ID, sPrepare.g.Emj["RIGHT"])
}

// newRoleEmbed 함수는 role guide와 현재 게임에 추가된 직업 / 게임의 참여중인 인원수 + 3 임베드를 만든다
func (sPrepare *Prepare) NewRoleEmbed() *embed.Embed {
	roleEmbed := embed.NewEmbed()
	roleEmbed.SetTitle("직업 추가")
	roleEmbed.AddField(sPrepare.g.RG[sPrepare.roleIndex].RoleName, strings.Join(sPrepare.g.RG[sPrepare.roleIndex].RoleGuide, "\n"))
	roleStr := ""
	if len(sPrepare.g.RoleView) == 0 {
		roleStr += "*추가된 직업이 없습니다.*"
	}
	for _, item := range sPrepare.g.RoleView {
		cnt := len(sPrepare.g.GetRoleUsers(item))
		roleStr += item.String() + " " + strconv.Itoa(cnt) + "개"
		if cnt == sPrepare.g.RG[sPrepare.roleIndex].Max {
			roleStr += " 최대"
		}
		roleStr += "\n"
	}
	roleEmbed.AddField("추가된 직업", roleStr)
	roleEmbed.SetFooter("현재 인원에 맞는 직업 수: " + strconv.Itoa(len(sPrepare.g.RoleView)) + " / " + strconv.Itoa(len(sPrepare.g.UserList)+3))
	return roleEmbed
}

// newEnterEmbed 함수는 게임 참여자 목록 임베드를 만든다
func (sPrepare *Prepare) NewEnterEmbed() *embed.Embed {
	enterEmbed := embed.NewEmbed()
	enterEmbed.SetTitle("게임 참가")
	enterEmbed.AddField("참가자 목록", "현재 참가 인원: "+strconv.Itoa(len(sPrepare.g.UserList))+"명\n")
	enterStr := ""
	if len(sPrepare.g.UserList) == 0 {
		enterStr += "*참가자가 없습니다.*"
	}
	for _, item := range sPrepare.g.UserList {
		enterStr += "`" + item.nick + "`\n"
	}
	enterEmbed.AddField("참가자 목록", enterStr)
	enterEmbed.SetFooter("⭕: 입장 ❌: 퇴장")
	return enterEmbed
}
