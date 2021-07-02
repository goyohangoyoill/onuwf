package game

// RoleSentinel is one of role in wfgame
type RoleSentinel struct {
	Role
}

// Action is role action
func (rs RoleSentinel) Action(tar TargetObject, player *User, g *Game) {
	//session 메세지는 state에서 보낼거임
	//action에서는 game 상태 바꾸는 action만
}

func (rs RoleSentinel) GenLog(tar TargetObject, player *User, g *Game) {
}

// String function that return role name in korean
func (rs RoleSentinel) String() string {
	return "수호자"
}
