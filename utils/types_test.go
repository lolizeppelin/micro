package utils

import (
	"fmt"
	"reflect"
	"testing"
)

func TestTypes(t *testing.T) {

	var a []string
	var b []byte

	x := reflect.TypeOf([]byte{})
	y := reflect.TypeOf(a)
	z := reflect.TypeOf(b)

	fmt.Printf("bytes check %v\n", TypeOfBytes == x)
	fmt.Printf("bytes check %v\n", TypeOfBytes == y)
	fmt.Printf("bytes check %v\n", TypeOfBytes == z)
	fmt.Printf("bytes check %v\n", TypeOfBytes == reflect.TypeOf([]byte{}))

}
