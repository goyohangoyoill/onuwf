module ONUWF

go 1.16

require (
	github.com/bwmarrin/discordgo v0.23.2 // indirect
	github.com/clinet/discordgo-embed v0.0.0-20190411043415-d754bc1a576c // indirect
	go.mongodb.org/mongo-driver v1.5.3 // indirect
	onuwf.com/game v0.0.0
	onuwf.com/util v0.0.0
)

replace (
	onuwf.com/game v0.0.0 => ./game
	onuwf.com/util v0.0.0 => ./util
)
