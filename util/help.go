package util

import (
	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"

	wfGame "onuwf.com/game"
)

const (
	prefix = "ㅁ"
)

func ReadJSON(rg []wfGame.RoleGuide) {
	// 명령어
	readCommandJSON()
	// 게임방법
	readRuleJSON()
	// 참고
	readNoteJSON(rg)
	// 게임배경
	readBackgroundJSON()
	// 승리조건
	readVictoryConditionJSON()
}

// 게임진행에 관련된 명령어가 입력될시 각 명령어에 해당하는 메시지를 출력하는 함수
func PrintHelpList(s *discordgo.Session, m *discordgo.MessageCreate, rg []wfGame.RoleGuide) bool {
	// "ㅁ직업소개 <직업명>" or "ㅁ직업소개 모두"
	if printRoleInfo(s, m, rg) {
		return true
	}
	switch m.Content {
	case prefix + "도움말":
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed(backgroundTitle, backgroundMsg))
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed(noteTitle, noteMsg))
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed(ruleTitle, ruleMsg))
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed(commandTitle, commandMsg))
	case prefix + "명령어":
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed(commandTitle, commandMsg))
	case prefix + "게임배경":
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed(backgroundTitle, backgroundMsg))
	case prefix + "게임방법":
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed(ruleTitle, ruleMsg))
	case prefix + "참고":
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed(noteTitle, noteMsg))
	case prefix + "승리조건":
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed(vcTitle, vcMsg))
	case prefix + "직업목록":
		printRoleList(s, m, rg)
	case prefix + "능력순서":
		printSkillOrder(s, m, rg)
	default:
		return false
	}
	return true
}