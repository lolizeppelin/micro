package log

import (
	"fmt"
	"github.com/lolizeppelin/micro/utils"
	"runtime"
	"strings"
)

const logrusPackage = "github.com/sirupsen/logrus"
const libsPackage = "go-diyibo/libs/libs/logging"

var skips = []string{
	logrusPackage,
	libsPackage,
}

func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}

	return f
}

func caller(*runtime.Frame) string {

	pcs := make([]uintptr, 20)
	depth := runtime.Callers(4, pcs)
	frames := runtime.CallersFrames(pcs[:depth])
	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)
		// If the caller isn't part of this package, we're done
		if !utils.IncludeInSlice(skips, pkg) {
			return fmt.Sprintf(" %s:%d ", f.File, f.Line)
		}
	}
	return ""
}
