package game

import (
	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
	"fmt"
)

// Prepare is test
type StateVote struct {
	g            *Game
	voted_list	[]int
	user_num	int
}

// PressNumBtn 사용자가 숫자 이모티콘을 눌렀을 때 state에서 하는 동작
func (v StateVote) PressNumBtn(s *discordgo.Session, r *discordgo.MessageReactionAdd, num int) {
	//num를 받음
	//해당 index list count +1
	v.voted_list[num]++
	fmt.Println(v.voted_list[num])
}

// PressDisBtn 사용자가 버려진 카드 이모티콘을 눌렀을 때 state에서 하는 동작
func (v StateVote) PressDisBtn(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	//do nothing
}


// PressYesBtn 사용자가 yes 이모티콘을 눌렀을 때 state에서 하는 동작
func (v StateVote) PressYesBtn(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	//do nothing
}

// PressNoBtn 사용자가 No 이모티콘을 눌렀을 때 state에서 하는 동작
func (v StateVote) PressNoBtn(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	//do nothing	
}


// PressDirBtn 좌 -1, 우 1 사용자가 좌우 방향 이모티콘을 눌렀을 때 state에서 하는 동작
func (v StateVote) PressDirBtn(s *discordgo.Session, r *discordgo.MessageReactionAdd, dir int) {
	//do nothing
}

func (v StateVote) VoteProcess(s *discordgo.Session, g *Game) {
	//send msg
	//개별로 각 채널에서 수행하게 해야함
	//user vote event handler
	//참가자마다 입력을 처리 *go routine 입력을 state vote preesnumbtn 이용해서 count
	// 결과 값을 visualization (통합채널) 

	v.voted_list = make([]int, v.user_num)
	voteEmbed := embed.NewEmbed()
	voteEmbed.SetTitle("투표")
	voteMsg, _ := s.ChannelMessageSendEmbed(g.ChanID, voteEmbed.MessageEmbed)
	//addNumAddEmoji(s, voteMsg)
	//fmt.Println(r.UserID, r.MessageID, r.ChannelID, r.GuildID)
}

func addNumAddEmoji(s *discordgo.Session, msg *discordgo.Message) {
	emj = EmojiInit()
	s.MessageReactionAdd(msg.ChannelID, msg.ID, emj["n1"])
	s.MessageReactionAdd(msg.ChannelID, msg.ID, emj["n2"])
	s.MessageReactionAdd(msg.ChannelID, msg.ID, emj["n3"])
	s.MessageReactionAdd(msg.ChannelID, msg.ID, emj["n4"])
}
