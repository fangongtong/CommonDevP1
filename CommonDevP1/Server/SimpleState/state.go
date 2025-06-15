package SimpleState

type IState interface {
	SwitchIn(context []interface{})
	SwitchOut()
}

type StateMgr struct {
	current IState
	next    IState
	context []interface{}
}

func NewStateMgr() *StateMgr {
	return &StateMgr{}
}

func (this *StateMgr) Run(first IState) {
	this.current = first
	this.current.SwitchIn(this.context)
	for {
		this.current.SwitchOut()
		this.current, this.next = this.next, nil
		this.current.SwitchIn(this.context)
	}
}

func (this *StateMgr) Switch(next IState, context []interface{}) {
	this.next = next
	this.context = context
}
