module ONUWF

go 1.16

require (
	github.com/bwmarrin/discordgo v0.23.2
	github.com/clinet/discordgo-embed v0.0.0-20190411043415-d754bc1a576c // indirect
	github.com/goyohangoyoill/ONUWF/game v0.0.0-20210716193736-9feb9e0ae0e4
	github.com/goyohangoyoill/ONUWF/util v0.0.0-20210717065002-b5ff528a2b28
	github.com/goyohangoyoill/ONUWF/util/json v0.0.0-00010101000000-000000000000
	github.com/mitchellh/mapstructure v1.4.1 // indirect
)

replace github.com/goyohangoyoill/ONUWF/game => ./game

replace github.com/goyohangoyoill/ONUWF/util => ./util

replace github.com/goyohangoyoill/ONUWF/util/json => ./util/json
