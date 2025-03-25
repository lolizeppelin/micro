package micro

import (
	"errors"
	"fmt"
	"github.com/lolizeppelin/micro/utils"
	"strings"
)

var (
	NoVersion  = errors.New("no version")
	VersionErr = errors.New("version value error")
)

type Version struct {
	Major int `json:"major,omitempty"`
	Minor int `json:"minor,omitempty"`
	Patch int `json:"patch,omitempty"`
}

func (v Version) Main() string {
	return utils.UnsafeToString(v.Minor)
}

// Version 版本号字符串,参数用于是否输出patch版本
func (v Version) Version(patch ...bool) string {
	if len(patch) > 0 && patch[0] {
		return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	}
	return fmt.Sprintf("%d.%d", v.Major, v.Minor)
}

func (v Version) Compare(version Version, patch ...bool) int {
	if v.Major > version.Major {
		return 1
	} else if v.Major < v.Major {
		return -1
	}
	if v.Minor > version.Minor {
		return 1
	} else if v.Minor < v.Minor {
		return -1
	}
	if len(patch) > 0 && patch[0] {
		if v.Patch > version.Patch {
			return 1
		} else if v.Patch < v.Patch {
			return -1
		}
	}
	return 0
}

// MarshalJSON Implementing the json.Marshaler interface
func (v Version) MarshalJSON() ([]byte, error) {
	// Use the standard library to marshal the map to JSON
	return []byte(fmt.Sprintf("{\"major\": %d, \"minor\": %d, \"patch\": %d}", v.Major, v.Minor, v.Patch)), nil
}

func NewVersion(v string) (*Version, error) {

	var (
		major int
		minor int
		patch int
		err   error
	)
	if v == "" {
		return nil, NoVersion
	}
	s := strings.Split(v, ".")
	parts := len(s)
	if len(s) <= 0 || parts > 3 {
		return nil, VersionErr
	}

	major, err = utils.StringToInt(s[0])
	if err != nil {
		return nil, VersionErr
	}
	if major <= 0 {
		return nil, VersionErr
	}
	if parts >= 2 {
		minor, err = utils.StringToInt(s[1])
		if err != nil {
			return nil, VersionErr
		}
		if minor < 0 {
			return nil, VersionErr
		}
	}
	if parts >= 3 {
		patch, err = utils.StringToInt(s[2])
		if err != nil {
			return nil, VersionErr
		}
		if patch < 0 {
			return nil, VersionErr
		}
	}

	return &Version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}, nil

}
