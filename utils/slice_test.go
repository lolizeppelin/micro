package utils

import (
	"fmt"
	"testing"
)

func TestSC(t *testing.T) {

	fmt.Println(NumberSliceJoin([]int64{
		1, 2, 3, 4, 10, 200,
	}))
}
