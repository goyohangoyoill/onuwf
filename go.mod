module ONUWF

go 1.16

require (
	github.com/bwmarrin/discordgo v0.23.2
	onuwf.com/game v0.0.0
	onuwf.com/util v0.0.0
)

replace (
	onuwf.com/game v0.0.0 => ./game
	onuwf.com/util v0.0.0 => ./util
)
