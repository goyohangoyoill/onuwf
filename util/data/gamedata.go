package data

import (
	"time"
)

type GameData struct {
	ObjectID           string
	StartTime, EndTime time.Time
	GuildID            string
	ChanID             string
	MasterID           string
	RoleList           []string
	UserList           []User
	OriDisRole         []string
	LastDisRole        []string
}

type User struct {
	UID      string
	Nick     string
	OriRole  string
	LastRole string
	isWin    bool
}

type UserData struct {
	UID            string
	Nick           string
	Title          string
	RecentGameTime time.Time
	CntPlay        int
	CntWin         int
	LastRoleList   []string
	PlayedGameOID  []string
}
