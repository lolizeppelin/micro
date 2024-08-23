package utils

import (
	"math"
	"time"
)

var _start = time.Now()

func Monotonic() time.Duration { // 单调递增时间
	return time.Now().Sub(_start)
}

func NowUnix() int64 {
	return time.Now().Unix()
}

func BackoffDelay(attempts int) time.Duration {
	if attempts > 13 {
		return 2 * time.Minute
	}
	return time.Duration(math.Pow(float64(attempts), math.E)) * time.Millisecond * 100
}

func ExpireTime(timeout int64) (int64, int64) {
	now := NowUnix()
	return now, now + timeout

}
