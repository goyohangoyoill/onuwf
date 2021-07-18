package game

// Hunter 는 한밤의 늑대인간 중 <사냥꾼> 에 대한 객체이다.
type Hunter struct {
	id int
}

// SendUserSelectGuide 직업 능력을 발휘하기 위한 선택지를 보내는 함수
func (h *Hunter) SendUserSelectGuide(player *User, g *Game, pageNum int) string {
	return "Hunter have no Msg"
}

// Action 함수는 <사냥꾼> 의 특수능력 사용에 대한 함수이다.
func (h *Hunter) Action(tar *TargetObject, player *User, g *Game) {
}

// GenLog 함수는 <사냥꾼> 의 특수능력 사용에 대한 함수이다.
func (h *Hunter) GenLog(tar *TargetObject, player *User, g *Game) {
}

// String 함수는 <사냥꾼> 문자열을 반환하는 함수이다.
func (h *Hunter) String() string {
	return "사냥꾼"
}

// ID 함수는 <사냥꾼> 의 고유값을 반환하는 함수이다.
func (h *Hunter) ID() int {
	return h.id
}
