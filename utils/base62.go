package utils

func ToBase62Bytes(n uint) (s []byte) {
	if n <= 0 {
		s = []byte("0")
		return
	}
	for n > 0 {
		s = append(s, baseChars[n%base])
		n = n / base
	}
	for i := 0; i < len(s)/2; i++ {
		s[i], s[len(s)-i-1] = s[len(s)-i-1], s[i]
	}
	return
}

func ToBase62(n uint) string {
	return string(ToBase62Bytes(n))
}
