package game

// User : 각 사용자들의 정보를 담고 있는 스트럭처
type User struct {
	// 각 유저의 UID
	userID string

	// 각 유저의 닉네임
	nick string

	// 각 유저가 속한 게임이 진행중인 채널 ID
	chanID string

	// 각 유저의 DM 채널 ID
	dmChanID string
}

// NewUser make new user object
func NewUser(uid, nick, chanID, dmChanID string) (u *User) {
	u.userID = uid
	u.nick = nick
	u.chanID = chanID
	u.dmChanID = dmChanID
	return
}
