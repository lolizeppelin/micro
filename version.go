package micro

import (
	"fmt"
	"github.com/lolizeppelin/micro/utils"
	"strings"
)

type Version struct {
	major int
	minor int
}

func (v *Version) Main() string {
	return utils.UnsafeToString(v.major)
}

// Major 主版本
func (v *Version) Major() int {
	return v.major
}

// Minor 次版本
func (v *Version) Minor() int {
	return v.minor
}

func (v *Version) Version() string {
	return fmt.Sprintf("%d.%d", v.major, v.minor)
}

func NewVersion(v string) (*Version, error) {

	var (
		major int
		minor int
		err   error
	)
	s := strings.Split(v, ".")
	if len(s) <= 0 || len(s) > 2 {
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

	return &Version{
		major: major,
		minor: minor,
	}, nil

}
