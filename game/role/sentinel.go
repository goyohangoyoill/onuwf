package game

type RoleSentinel struct {
	*Role
}

func (r *RoleSentinel) Action(tar TargetObject, player *User, g *Game) {
	//session 메세지는 state에서 보낼거임
	//action에서는 game 상태 바꾸는 action만
}

func (r *RoleSentinel) String() string {
	return "수호자"
}
