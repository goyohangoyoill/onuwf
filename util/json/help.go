// Package json is a package for json files
package json

import (
	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
)

// ReadJSON ./asset에 있는 json파일들을 읽어오는 함수
func ReadJSON(rg []RoleGuide, prefix string) {
	// 명령어
	readCommandJSON(prefix)
	// 게임방법
	readRuleJSON(prefix)
	// 참고
	readNoteJSON(rg)
	// 게임배경
	readBackgroundJSON(prefix)
	// 승리조건
	readVictoryConditionJSON()
	// 나무위키
	readWikiJSON()
}

// PrintHelpList 게임진행에 관련된 명령어가 입력될시 각 명령어에 해당하는 메시지를 출력하는 함수
func PrintHelpList(s *discordgo.Session, m *discordgo.InteractionCreate, rg []RoleGuide, prefix string) bool {
	// "ㅁ직업소개 <직업명>" or "ㅁ직업소개 모두"
	if printRoleInfo(s, m, rg, prefix) {
		return true
	}
	switch m.ApplicationCommandData().Name {
	case "도움말":
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed(backgroundTitle, backgroundMsg))
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed(noteTitle, noteMsg))
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed(ruleTitle, ruleMsg))
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed(commandTitle, commandMsg))
	case "명령어":
		fallthrough
	case "help":
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed(commandTitle, commandMsg))
	case "게임배경":
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed(backgroundTitle, backgroundMsg))
	case "게임방법":
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed(ruleTitle, ruleMsg))
	case "참고":
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed(noteTitle, noteMsg))
	case "승리조건":
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed(vcTitle, vcMsg))
	case "직업목록":
		printRoleList(s, m, rg, prefix)
	case "직업순서", "순서", "능력순서":
		printSkillOrder(s, m, rg, prefix, false)
	case "직업서순", "서순", "능력서순":
		printSkillOrder(s, m, rg, prefix, true)
	case "나무위키":
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbed(wikiTitle, wikiMsg))
	default:
		return false
	}
	return true
}
