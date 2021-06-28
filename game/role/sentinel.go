// +build linux,amd64,go1.15,!cgo

package role

type RoleSentinel struct {
	role
}

func (r *RoleSentinel) Action(tar targetObject, player *user, g *game) {
	//session 메세지는 state에서 보낼거임
	//action에서는 game 상태 바꾸는 action만
}

func (r *RoleSentinel) String() string {
	return "수호자"
}
