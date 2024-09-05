package micro

import (
	"encoding/json"
	"fmt"
	"github.com/lolizeppelin/micro/utils"
	"strings"
)

type Version struct {
	major int
	minor int
	patch int
}

func (v Version) Main() string {
	return utils.UnsafeToString(v.major)
}

// Major 主版本
func (v Version) Major() int {
	return v.major
}

// Minor 次版本
func (v Version) Minor() int {
	return v.minor
}

// Patch 次版本
func (v Version) Patch() int {
	return v.patch
}

// Version 版本号字符串,参数用于是否输出patch版本
func (v Version) Version(patch ...bool) string {
	if len(patch) > 0 && patch[0] {
		return fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)
	}
	return fmt.Sprintf("%d.%d", v.major, v.minor)
}

func (v Version) Compare(version Version) int {
	if v.major > version.major {
		return 1
	} else if v.major < v.major {
		return -1
	}
	if v.minor > version.minor {
		return 1
	} else if v.minor < v.minor {
		return -1
	}
	if v.patch > version.patch {
		return 1
	} else if v.patch < v.patch {
		return -1
	}
	return 0
}

// MarshalJSON Implementing the json.Marshaler interface
func (v Version) MarshalJSON() ([]byte, error) {
	// Create a map to hold the JSON representation
	m := map[string]int{
		"major": v.major,
		"minor": v.minor,
		"patch": v.patch,
	}
	// Use the standard library to marshal the map to JSON
	return json.Marshal(m)
}

func NewVersion(v string) (*Version, error) {

	var (
		major int
		minor int
		patch int
		err   error
	)
	s := strings.Split(v, ".")
	if len(s) <= 0 || len(s) > 3 {
		return nil, fmt.Errorf("version value error")
	}

	major, err = utils.StringToInt(s[0])
	if err != nil {
		return nil, fmt.Errorf("major version error")
	}
	if major <= 0 {
		return nil, fmt.Errorf("major value lt 0")
	}
	if len(s) >= 2 {
		minor, err = utils.StringToInt(s[1])
		if err != nil {
			return nil, fmt.Errorf("minor version error")
		}
		if minor < 0 {
			return nil, fmt.Errorf("minor value less then 0")
		}
	}
	if len(s) >= 3 {
		patch, err = utils.StringToInt(s[2])
		if err != nil {
			return nil, fmt.Errorf("patch version error")
		}
		if patch < 0 {
			return nil, fmt.Errorf("patch value less then 0")
		}
	}

	return &Version{
		major: major,
		minor: minor,
	}, nil

}
