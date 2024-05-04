package utils

import (
	"errors"
	"os"
	"regexp"
	"strings"
	"unicode"
)

var _unixPathLikeRegex = regexp.MustCompile(`^/([^\p{C}\p{Z}]*[\p{L}\p{N}\.\-]+[^\p{C}\p{Z}]*)+$`)

func PathExist(path string) (os.FileInfo, error) {
	info, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return info, nil
}

func PathIsDir(path string) (bool, error) {
	fileInfo, err := PathExist(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	if fileInfo == nil {
		return false, nil
	}
	return fileInfo.IsDir(), errors.New("path is not folder")
}

func PathIsRegularFile(path string) (bool, error) {
	fileInfo, err := PathExist(path)
	if err != nil {
		return false, err
	}
	if fileInfo == nil {
		return false, nil
	}
	return fileInfo.Mode().IsRegular(), nil
}

func PathIsUnixSockFile(path string) (bool, error) {
	fileInfo, err := os.Lstat(path)
	if err != nil {
		return false, err
	}
	return isUnixSocketFile(fileInfo), nil
}

func MakeDirs(path string) error {
	return os.MkdirAll(path, 0755)
}

func IsSafeUnixLikePath(path string) bool {
	if len(path) <= 0 {
		return false
	}
	if strings.HasSuffix(path, "/") {
		return false
	}
	if !_unixPathLikeRegex.MatchString(path) {
		return false
	}
	isPrevSlash := false
	isFirstChar := true
	var stringBuffer strings.Builder

	for _, r := range path {
		if !unicode.IsPrint(r) {
			return false
		}
		if r == '/' {
			if isPrevSlash {
				return false
			}
			isPrevSlash = true
			isFirstChar = true
			if stringBuffer.Len() > 0 {
				if stringBuffer.String() == "." || stringBuffer.String() == ".." {
					return false
				}
				bufferStr := stringBuffer.String()
				if bufferStr[len(bufferStr)-1] == '-' {
					return false
				}
				stringBuffer.Reset()
			}
		} else {
			isPrevSlash = false
			if r == '-' && isFirstChar {
				return false
			}
			isFirstChar = false
			stringBuffer.WriteRune(r)
		}
	}
	if stringBuffer.Len() > 0 {
		bufferStr := stringBuffer.String()
		if bufferStr[len(bufferStr)-1] == '-' {
			return false
		}
	}

	return true
}

func ReadFile(path string) ([]byte, error) {
	ok, _ := PathIsRegularFile(path)
	if !ok {
		return nil, errors.New("path is not regular file")
	}
	return os.ReadFile(path)
}
