module game

go 1.16

require (
	github.com/bwmarrin/discordgo v0.23.2
	github.com/clinet/discordgo-embed v0.0.0-20190411043415-d754bc1a576c
	github.com/goyohangoyoill/onuwf/util v0.0.0-20210804134634-b8cd3cfb31a7
	github.com/goyohangoyoill/onuwf/util/json v0.0.0-00010101000000-000000000000
	github.com/mitchellh/mapstructure v1.4.1 // indirect
)

replace github.com/goyohangoyoill/onuwf/util => ../util

replace github.com/goyohangoyoill/onuwf/util/json => ../util/json
