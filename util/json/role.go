// Package json is a package for json files
package json

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
)

// RoleGuide has info of each role
type RoleGuide struct {
	RoleName  string   `json:"roleName"`
	RoleGuide []string `json:"roleGuide"`
	Max       int      `json:"max"`
	Faction   string   `json:"faction"`
	Priority  int      `json:"priority"`
}

// RoleGuideInit 직업 가이드 에셋 불러오기.
func RoleGuideInit(rg *[]RoleGuide) {
	rgFile, err := os.Open("asset/role_guide.json")
	if err != nil {
		log.Fatal(err)
	}
	defer rgFile.Close()
	var byteValue []byte
	byteValue, err = ioutil.ReadAll(rgFile)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal([]byte(byteValue), rg)
}

// RoleName을 이용하여 RoleGuide를 가져옴
func roleGuide(role string, rg []RoleGuide) string {
	guide := ""
	for i := 0; i < len(rg); i++ {
		if rg[i].RoleName == role {
			for j := 0; j < len(rg[i].RoleGuide); j++ {
				guide += rg[i].RoleGuide[j]
				if j-1 < len(rg[i].RoleGuide) {
					guide += "\n"
				}
			}
		}
	}
	return guide
}

// 모든 RoleName을 능력순서대로 가져옴
func roleList(rg []RoleGuide) []string {
	var list []string
	for i := 0; i < len(rg); i++ {
		if rg[i].Priority != -1 {
			list = append(list, rg[i].RoleName)
		}
	}
	return list
}

// "ㅁ직업소개 <직업명>", "ㅁ직업소개" 입력시 실행되는 함수
// 메세지 출력시: true, 미출력시: false
func printRoleInfo(s *discordgo.Session, m *discordgo.InteractionCreate, rg []RoleGuide, prefix string) bool {
	// "ㅁ직업소개"만 입력시
	classStr := strings.Split(m.ApplicationCommandData().Name, " ")
	if len(classStr) == 1 {
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed("직업소개", prefix+"직업소개 <직업명> 으로 요청하세요."))
		return true
	}
	classList := roleList(rg)
	// ㅁ직업소개 모두
	if classStr[1] == "모두" {
		uChan, _ := s.UserChannelCreate(m.Member.User.ID)
		for _, item := range classList {
			s.ChannelMessageSendEmbed(uChan.ID, embed.NewGenericEmbed("**"+item+"**", roleGuide(item, rg)))
		}
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed("모든 직업 소개가 DM으로 전송되었습니다.", ""))
		return true
	}
	// ㅁ직업소개 <직업명>
	for _, item := range classList {
		if classStr[1] == item {
			s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed("**"+item+" 소개**", roleGuide(item, rg)))
			return true
		}
	}
	// 직업명 잘못 입력시
	s.ChannelMessageSend(m.ChannelID, "직업 이름이 잘못되었습니다.")
	return true
}

// "ㅁ직업목록" 명령어 입력시 실행되는 함수
func printRoleList(s *discordgo.Session, m *discordgo.InteractionCreate, rg []RoleGuide, prefix string) {
	roleList := roleList(rg)
	printMsg := ""
	for i, item := range roleList {
		printMsg += item + " "
		if i%4 == 3 {
			printMsg += "\n"
		}
	}
	printMsg += "\n\n`" +
		prefix + "직업소개 <직업명>` 으로 직업소개 불러오기" +
		"\n`" + prefix + "직업소개 모두` 로 모든 직업소개 DM으로 받기"
	s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed("**구현된 직업 목록**", printMsg))
}

// 능력순서 임베드를 출력하는 함수
func printSkillOrder(s *discordgo.Session, m *discordgo.InteractionCreate, rg []RoleGuide, prefix string, isReversed bool) {
	printTitle := ""
	printMsg := ""
	roleList := roleList(rg)
	if isReversed {
		printTitle = "특수능력 사용 서순"
		printMsg += "투표시작"
		for i := len(roleList) - 1; i >= 0; i-- {
			if rg[i].Priority != 4242 {
				printMsg += " <- " + roleList[i]
				if i%3 == 2 {
					printMsg += "\n"
				}
			}
		}
	} else {
		printTitle = "특수능력 사용 순서"
		for i, item := range roleList {
			if rg[i].Priority != 4242 {
				printMsg += item + " -> "
				if i%3 == 2 {
					printMsg += "\n"
				}
			}
		}
		printMsg += "투표시작"
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed("**"+printTitle+"**", printMsg))
}
