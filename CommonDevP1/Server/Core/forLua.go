package Core

import (
	"time"

	"github.com/yuin/gopher-lua"
)

// api2L_ 表示开放给lua调用的go函数

func Api2L_Sleep(L *lua.LState) int {
	v := L.ToInt(1)
	time.Sleep(time.Millisecond * time.Duration(v))
	return 0
}
