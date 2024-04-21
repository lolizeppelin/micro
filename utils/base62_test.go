package utils

import (
	"fmt"
	"testing"
)

type TestValue struct {
	value int
}

func warp(key string, value *TestValue) func() {

	v := value
	return func() {
		fmt.Printf("key %s, value %d\n", key, v.value)
	}
}

func TestB62(t *testing.T) {

	fmt.Printf("~~~ %s, %d \n", ToBase62(1), 1)
	fmt.Printf("~~~ %s, %d \n", ToBase62(20), 20)

	value, err := FromBase62("001")
	if err != nil {
		fmt.Printf("from base 62 failed %s", err.Error())
		return
	}
	fmt.Printf("~~~ %d, 001 \n", value)

	m := map[string]TestValue{
		"a": {value: 1},
		"b": {value: 2},
		"c": {value: 3},
	}
	var fs []func()

	var values []*TestValue

	for k, v := range m {

		vv := v

		values = append(values, &vv)
		fn := warp(k, &vv)
		fs = append(fs, fn)

	}

	for _, fn := range fs {
		fn()
	}

	fmt.Println("--------")

	for _, v := range values {
		fmt.Printf("vlaue %d\n", v.value)
	}

	sl := []string{"a", "b", "c"}
	il := []int{1, 2, 3}
SKIP:
	for _, s := range sl {
		fmt.Printf("s ~~%s\n", s)
		for _, i := range il {
			fmt.Printf("i ~~%d\n", i)
			if i == 2 {
				continue SKIP
			}
		}

	}

}
