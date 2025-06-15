package Core

import (
	"time"
)

type ILoop interface {
	Loop()
}

type _Looper struct {
	runFlg bool
	looper []ILoop
}

func (this *_Looper) AddLooper(l ILoop) {
	this.looper = append(this.looper, l)
}

func (this *_Looper) Run() {
	this.runFlg = true

	for this.runFlg {
		for _, l := range this.looper {
			this.looper.Loop()
		}
	}
}

func (this *_Looper) Stop() {
	this.runFlg = false
}
