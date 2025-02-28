package utils

import (
	"math"
	"time"
)

const (
	TimestampFormat = "2006-01-02 15:04:05"
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

func StringToTime(s string) (time.Time, error) {
	return time.Parse(TimestampFormat, s)
}

// GetDay 获取当日零点
func GetDay(at ...int64) int64 {
	var t time.Time
	if len(at) > 0 && at[0] > 0 {
		t = time.Unix(at[0], 0)
	} else {
		t = time.Now()
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix()
}

// GetMonday 获取本周一
func GetMonday(at ...int64) int64 {
	var pos time.Time
	if len(at) > 0 && at[0] > 0 {
		pos = time.Unix(at[0], 0)
	} else {
		pos = time.Now()
	}
	// 计算本周的周一
	// time.Weekday 中，周一到周日分别对应 1 到 7，但 time.Monday 的值为 1
	weekday := pos.Weekday()             // 当前是周几
	passed := int(weekday - time.Monday) // 距离周一有多少天
	if passed < 0 {
		passed += 7 // 如果当前是周日，需要调整
	}
	monday := pos.AddDate(0, 0, -passed) // 减去天数得到周一
	// 将时间部分设置为 00:00:00
	return time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, monday.Location()).Unix()
}

// GetMonth 获取月初
func GetMonth(at ...int64) int64 {
	// 获取本月的月初
	if len(at) > 0 && at[0] > 0 {
		t := time.Unix(at[0], 0)
		return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location()).Unix()
	}
	t := time.Now()
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location()).Unix()
}

// GetYear 获取今年第一天
func GetYear(at ...int64) int64 {
	// 今年第一天
	if len(at) > 0 && at[0] > 0 {
		t := time.Unix(at[0], 0)
		return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location()).Unix()
	}
	t := time.Now()
	return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location()).Unix()

}

// GetNextDay 获取明日零点
func GetNextDay(at ...int64) int64 {
	// 获取下一天的开始时间
	return GetDay(at...) + 86400
}

// GetNextYear 获取次年
func GetNextYear(at ...int64) int64 {
	var t time.Time
	if len(at) > 0 && at[0] > 0 {
		t = time.Unix(at[0], 0)
	} else {
		t = time.Now()
	}
	// 获取下一年的年初
	return time.Date(t.Year()+1, 1, 1, 0, 0, 0, 0, t.Location()).Unix()
}

// GetNextMonth 获取下个月
func GetNextMonth(at ...int64) int64 {
	var t time.Time
	if len(at) > 0 && at[0] > 0 {
		t = time.Unix(at[0], 0)
	} else {
		t = time.Now()
	}

	// 处理月份溢出
	year := t.Year()
	month := t.Month() + 1
	if month > 12 {
		month = 1
		year++
	}
	// 获取下个月的月初
	return time.Date(year, month, 1, 0, 0, 0, 0, t.Location()).Unix()
}
