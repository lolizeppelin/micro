package utils

import (
	"fmt"
	"testing"
)

func TestQuery(t *testing.T) {

	m, err := ParseQuery("a=2&b=3&c0=4")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for k, v := range m {
		fmt.Printf(" key %s, value %s\n", k, v)
	}

}
