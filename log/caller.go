package log

import (
	"fmt"
	"golang.org/x/exp/slices"
	"runtime"
	"strings"
)

const (
	logrusPackage = "github.com/sirupsen/logrus"
	libsPackage   = "github.com/lolizeppelin/micro/log"
)

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
		if slices.Contains(skips, pkg) {
			continue
		}
		return fmt.Sprintf("%s:%d ", f.File, f.Line)
	}
	return ""
}
