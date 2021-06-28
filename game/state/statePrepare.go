// +build linux,amd64,go1.15,!cgo

package state

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
)

type StatePrepare struct {
	// state에서 가지고 있는 game
	g *game

	// factory 에서 쓰이게 될 role index
	roleIndex int

	// 직업추가 확인용 메세지
	roleAddMsg *discordgo.Message

	// 게임입장 확인용 메세지
	enterGameMsg *discordgo.Message
}

// PressNumBtn 사용자가 숫자 이모티콘을 눌렀을 때 StatePrepare에서 하는 동작
func (sPrepare *StatePrepare) PressNumBtn(s *discordgo.Session, r *discordgo.MessageReactionAdd, num int) {
	// do nothing
}

// PressDisBtn 사용자가 버려진 카드 이모티콘을 눌렀을 때 StatePrepare에서 하는 동작
func (sPrepare *StatePrepare) PressDisBtn(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// do nothing
}

// PressYesBtn 사용자가 yes 이모티콘을 눌렀을 때 StatePrepare에서 하는 동작
func (sPrepare *StatePrepare) PressYesBtn(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// 입장 메세지에서 리액션한거라면
	if r.MessageID == sPrepare.g.enterGameMsgID {
		// userList에 없으면 user append()
		if findUserIdx(r.UserID, sPrepare.g.userList) == -1 {
			//user 생성해서 append()
			userNick, _ := s.User(r.UserID)
			userDM, _ := s.UserChannelCreate(r.UserID)
			u := user{userID: r.UserID, nick: userNick.Username, chanID: r.ChannelID, dmChanID: userDM.ID}
			sPrepare.g.userList = append(sPrepare.g.userList, u)
			// 입장 확인 메세지 반영
			s.ChannelMessageEditEmbed(sPrepare.g.chanID, sPrepare.enterGameMsg.ID, newEnterEmbed(sPrepare.g).MessageEmbed)
		}
		// 직업추가 메세지에서 리액션한거라면
	} else if r.MessageID == sPrepare.g.roleAddMsgID {
		// roleFactory에서 현재 roleindex 위치 값을 받아 role 생성
		roleToAdd := rf.generateRole(sPrepare.roleIndex)
		// 추가된 role 개수가 max보다 작을 때만 추가
		if roleCount(roleToAdd, sPrepare.g.roleView) < rg[sPrepare.roleIndex].Max {
			// roleView는 ununique sorted니까 roleView에 중복된 상태로 sort index 찾아서 삽입
			for i, item := range sPrepare.g.roleView {
				if item.String() <= roleToAdd.String() {
					tmp := append(sPrepare.g.roleView[:i], roleToAdd)
					sPrepare.g.roleView = append(tmp, sPrepare.g.roleView[i:]...)
					break
				}
			}
			// roleSeq는 unique unsorted니까 roleSeq에 없으면 append
			if findRoleIdx(roleToAdd, sPrepare.g.roleSeq) == -1 {
				sPrepare.g.roleSeq = append(sPrepare.g.roleSeq, roleToAdd)
			}
			// 직업 추가 메세지 반영
			s.ChannelMessageEditEmbed(sPrepare.g.chanID, sPrepare.roleAddMsg.ID, newRoleEmbed(sPrepare.roleIndex, sPrepare.g).MessageEmbed)
		}
	}
}

// PressNoBtn 사용자가 No 이모티콘을 눌렀을 때 StatePrepare에서 하는 동작
func (sPrepare *StatePrepare) PressNoBtn(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// 입장 메세지에서 리액션한거라면
	if r.MessageID == sPrepare.g.enterGameMsgID {
		// userList에 있으면 지우기
		if i := findUserIdx(r.UserID, sPrepare.g.userList); i != -1 {
			// userList에서 지우고
			sPrepare.g.userList = append(sPrepare.g.userList[:i], sPrepare.g.userList[i+1:]...)
			// 입장 확인 메세지에서 지우기
			var msg string
			msg = sPrepare.enterGameMsg.Embeds[0].Description
			strings.Replace(msg, r.UserID+" 입장\n", "", 1)
			sPrepare.enterGameMsg.Embeds[0].Description = msg
			s.ChannelMessageEditEmbed(sPrepare.g.chanID, sPrepare.enterGameMsg.ID, sPrepare.enterGameMsg.Embeds[0])
		}
		// 직업추가 메세지에서 리액션한거라면
	} else if r.MessageID == sPrepare.g.roleAddMsgID {
		// roleFactory에서 현재 roleindex 위치 값을 받아 role 생성
		roleToRemove := rf.generateRole(sPrepare.roleIndex)
		// roleView는 ununique sorted니까 첫번째로 나오는거 찾아서 지우기
		if i := findRoleIdx(roleToRemove, sPrepare.g.roleView); i != -1 {
			sPrepare.g.roleView = append(sPrepare.g.roleView[:i], sPrepare.g.roleView[i+1:]...)
			// roleSeq는 unique unsorted니까 방금 지운 roleView에 없으면 지우기
		} else if i := findRoleIdx(roleToRemove, sPrepare.g.roleSeq); i != -1 {
			sPrepare.g.roleSeq = append(sPrepare.g.roleSeq[:i], sPrepare.g.roleSeq[i+1:]...)
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
				if num == rg[sPrepare.roleIndex].Max {
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
			s.ChannelMessageEditEmbed(sPrepare.g.chanID, sPrepare.roleAddMsg.ID, sPrepare.roleAddMsg.Embeds[0])
		}
	}
}

// PressDirBtn 좌 -1, 우 1 사용자가 좌우 방향 이모티콘을 눌렀을 때 StatePrepare에서 하는 동작
func (sPrepare *StatePrepare) PressDirBtn(s *discordgo.Session, r *discordgo.MessageReactionAdd, dir int) {
	if r.MessageID == sPrepare.g.enterGameMsgID {
		// 게임 시작
		if dir == 1 {
			if len(sPrepare.g.roleSeq) == len(sPrepare.g.userList)+3 {
				sPrepare.g.curState = &StatePlayable{g: sPrepare.g}
				s.ChannelMessageSendEmbed(sPrepare.g.chanID, embed.NewGenericEmbed("게임시작", ""))
			} else {
				s.ChannelMessageSendEmbed(sPrepare.g.chanID, embed.NewGenericEmbed("인원수 안 맞음", ""))
			}
		}
	} else if r.MessageID == sPrepare.g.roleAddMsgID {
		// roleindex 증감
		sPrepare.roleIndex += dir
		if dir > len(rg) {
			dir = 0
		} else if dir <= 0 {
			dir = len(rg)
		}
		s.ChannelMessageEditEmbed(sPrepare.g.chanID, sPrepare.g.roleAddMsgID,
			embed.NewGenericEmbed(rg[sPrepare.roleIndex].RoleName, strings.Join(rg[sPrepare.roleIndex].RoleGuide, "\n")))
	}
}
