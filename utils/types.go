package utils

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
