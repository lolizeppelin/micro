package utils

import (
	"bytes"
	"fmt"
)

// SizeBuff 有最大值限制的buff
type SizeBuff struct {
	max  int
	size int
	buff *bytes.Buffer
}

func (s *SizeBuff) Write(p []byte) (int, error) {
	n, err := s.buff.Write(p)
	if err != nil {
		return n, err
	}
	s.size += n
	if s.max > 0 && s.size > s.max {
		return n, fmt.Errorf("over size")
	}
	return n, err
}

func (s *SizeBuff) Bytes() []byte {
	return s.buff.Bytes()
}

func (s *SizeBuff) Reader() *bytes.Reader {
	return bytes.NewReader(s.Bytes())
}

func NewSizeBuff(size int) *SizeBuff {
	return &SizeBuff{
		max:  size,
		buff: new(bytes.Buffer),
	}
}
