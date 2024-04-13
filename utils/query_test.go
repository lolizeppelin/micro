package utils

import (
	"fmt"
	"testing"
)

func TestQuery(t *testing.T) {

	s := BuildQuery(map[string]string{
		"a": "=1",
	})
	fmt.Println(s)

}
