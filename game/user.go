package game

// User : 각 사용자들의 정보를 담고 있는 스트럭처
type User struct {
	// 각 유저의 UID
	UserID string
	// 각 유저의 닉네임
	nick string
	// 각 유저의 칭호
	title string
	// 각 유저가 속한 게임이 진행중인 채널 ID
	chanID string
	// 각 유저의 DM 채널 ID
	dmChanID string
	// 각 유저가 투표한 유저의 ID
	voteUserId string
}

func (u User) Nick() string {
	return u.nick
}

func (u User) Title() string {
	return u.title
}

// NewUser make new user object
func NewUser(uid, nick, chanID, dmChanID string) (u *User) {
	u = &User{}
	u.UserID = uid
	u.nick = nick
	u.chanID = chanID
	u.dmChanID = dmChanID
	return
}

func UpdateUser(cu *User, nick, title string) (u *User) {
	cu.nick = nick
	cu.title = title
	return cu
}
