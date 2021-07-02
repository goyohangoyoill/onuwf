package game

// Sentinel is one of role in wfgame
type Sentinel struct {
	id int
}

// Action is role action
func (r *Sentinel) Action(tar *TargetObject, player *User, g *Game) {
	//session 메세지는 state에서 보낼거임
	//action에서는 game 상태 바꾸는 action만
}

func (r *Sentinel) GenLog(tar *TargetObject, player *User, g *Game) {
}

// String function that return role name in korean
func (r *Sentinel) String() string {
	return "수호자"
}

// ID 함수는 <수호자> 의 고유값을 반환하는 함수이다.
func (r *Sentinel) ID() int {
	return r.id
}
