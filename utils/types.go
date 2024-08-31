package utils

import (
	"context"
	"github.com/golang/protobuf/proto"
	"reflect"
)

type IntType interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type UnsignedIntType interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type NumberType interface {
	~int | ~int8 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64
}

type SortAbleType interface {
	~int | ~int8 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64 | ~string
}

type BaseType interface {
	~int | ~int8 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64 | ~bool | ~string | ~complex64 | ~complex128
}

var (
	TypeOfError    = reflect.TypeOf((*error)(nil)).Elem()
	TypeOfBytes    = reflect.TypeOf(([]byte)(nil))
	TypeOfContext  = reflect.TypeOf(new(context.Context)).Elem()
	TypeOfProtoMsg = reflect.TypeOf(new(proto.Message)).Elem()

	TypesOfBaseType = []reflect.Kind{
		reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.String,
	}
)
