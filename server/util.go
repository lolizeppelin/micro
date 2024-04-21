package server

import (
	"fmt"
	"github.com/lolizeppelin/micro/utils"
)

const (
	MaxServerSN = 62*62*62 - 1
)

// SNBase62 62进制转换服务器id
func SNBase62[T utils.IntType](id T) string {
	sn := uint(id)
	if sn < 1 {
		sn = 1
	}
	if sn >= MaxServerSN {
		sn = MaxServerSN
	}
	sid := utils.ToBase62(sn)
	return fmt.Sprintf("%0*s", 3, sid)
}
