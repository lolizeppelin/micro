package utils

func Abs[T NumberType](i T) T {
	if i < 0 {
		return -i
	}
	return i
}
