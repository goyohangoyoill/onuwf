package game

// Villager 는 한밤의 늑대인간 중 <마을주민> 에 대한 객체이다.
type Villager struct {
	id int
}

// SendUserSelectGuide 직업 능력을 발휘하기 위한 선택지를 보내는 함수
func (v *Villager) SendUserSelectGuide(player *User, g *Game, pageNum int) string {
	return "Villager have no Msg"
}

// Action 함수는 <마을주민> 의 특수능력 사용에 대한 함수이다.
func (v *Villager) Action(tar *TargetObject, player *User, g *Game) {
	// Do Nothing
}

// GenLog 함수는 <마을주민> 의 특수능력 사용에 대한 함수이다.
func (v *Villager) GenLog(tar *TargetObject, player *User, g *Game) {
	// Do Nothing
}

// String 함수는 <마을주민> 문자열을 반환하는 함수이다.
func (v *Villager) String() string {
	return "마을주민"
}

// ID 함수는 <마을주민> 의 고유값을 반환하는 함수이다.
func (v *Villager) ID() int {
	return v.id
}
