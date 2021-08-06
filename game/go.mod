module game

go 1.16

require (
	github.com/bwmarrin/discordgo v0.23.2
	github.com/clinet/discordgo-embed v0.0.0-20190411043415-d754bc1a576c
	github.com/goyohangoyoill/ONUWF/util v0.0.0-20210804134634-b8cd3cfb31a7 // indirect
	github.com/goyohangoyoill/ONUWF/util/json v0.0.0-00010101000000-000000000000
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	golang.org/x/crypto v0.0.0-20200302210943-78000ba7a073 // indirect
	golang.org/x/sys v0.0.0-20190531175056-4c3a928424d2 // indirect
)

replace github.com/goyohangoyoill/ONUWF/util => ../util

replace github.com/goyohangoyoill/ONUWF/util/json => ../util/json
