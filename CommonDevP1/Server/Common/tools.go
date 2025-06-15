package Common

import (
	"fmt"
	"sync/atomic"
	"time"
)

var _UNIX_START = time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local).Unix()
var _ID_CNT uint32 = 0

func GetUID() string {
	now := time.Now().Unix() - _UNIX_START
	return fmt.Sprintf("%012X%04X", uint64(now), uint16(atomic.AddUint32(&_ID_CNT, 1)))
}
