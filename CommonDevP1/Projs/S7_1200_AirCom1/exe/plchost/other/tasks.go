package other

import (
	"fmt"
)

type TaskList struct {
	taskBaseInfoBuf []*TaskBaseInfo
}

func (this *TaskList) Init() {
	this.taskBaseInfoBuf = []*TaskBaseInfo{}
}

//  state 2:pause 3:run 4:over
func (this *TaskList) Notify(uuid string, state int) {

	if state == 4 {
		var tmp []*TaskBaseInfo = []*TaskBaseInfo{}
		for _, v := range this.taskBaseInfoBuf {
			if v.Uid != uuid {
				tmp = append(tmp, v)
			}
		}

		this.taskBaseInfoBuf = tmp
	} else {
		found := false
		for _, v := range this.taskBaseInfoBuf {
			if v.Uid == uuid {
				v.Oper = state
				found = true
				break
			}
		}
		if !found {
			this.taskBaseInfoBuf = append(this.taskBaseInfoBuf, &TaskBaseInfo{
				Uid:  uuid,
				Oper: state,
			})
		}
	}
	fmt.Println(uuid, state)
	fmt.Println("taskBaseInfoBuf len: ", len(this.taskBaseInfoBuf))
}

type TaskBaseInfo struct {
	Uid  string
	Oper int
}

func (this *TaskList) Tasks() []*TaskBaseInfo {
	return this.taskBaseInfoBuf
}
