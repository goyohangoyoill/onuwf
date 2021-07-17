module ONUWF

go 1.16

require (
	github.com/bwmarrin/discordgo v0.23.2
	github.com/goyohangoyoill/ONUWF/game v0.0.0-20210716193736-9feb9e0ae0e4
	github.com/goyohangoyoill/ONUWF/util v0.0.0-20210717065002-b5ff528a2b28
)

replace github.com/goyohangoyoill/ONUWF/game => ./game

replace github.com/goyohangoyoill/ONUWF/util => ./util
