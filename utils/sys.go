package utils

import "runtime"

const Linux = runtime.GOOS == "linux"
const Windows = runtime.GOOS == "windows"
