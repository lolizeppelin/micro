package utils

import (
	"github.com/sqids/sqids-go"
	"sync/atomic"
)

type generator struct {
	sequence *uint64
}

func (g *generator) next() uint64 {
	for {
		current := atomic.LoadUint64(g.sequence)
		next := current + 1

		// Reset to 0 if the next value exceeds 255
		if next > 255 {
			next = 0
		}

		// Attempt to update the sequence to the new value
		if atomic.CompareAndSwapUint64(g.sequence, current, next) {
			return next
		}
	}
}

type Sequence struct {
	builder   *sqids.Sqids
	generator map[uint64]*generator
	len       uint64
	counter   *uint64
}

func (s *Sequence) Next(a uint64, b uint64) (string, error) {
	n := atomic.AddUint64(s.counter, 1)
	index := n % s.len
	if index >= 65535 {
		atomic.StoreUint64(s.counter, 0)
	}
	g := s.generator[index]
	return s.builder.Encode([]uint64{
		a, b, g.next(), index,
	})
}

func (s *Sequence) Decode(id string) []uint64 {
	return s.builder.Decode(id)
}

/*
NewSequence 创建数字转字符串的序列化对象
@param chars 可用的字符串
@param length 输出最大长度
@param bulk 并发队列书香,默认4
*/
func NewSequence(chars string, length uint8, bulk ...uint16) (*Sequence, error) {
	builder, err := sqids.New(sqids.Options{
		Alphabet:  chars,
		MinLength: length,
	})
	if err != nil {
		return nil, err
	}
	bk := 4
	if len(bulk) > 0 && bulk[0] > 0 {
		bk = int(bulk[0])
	}
	s := &Sequence{
		builder:   builder,
		generator: map[uint64]*generator{},
		len:       uint64(bk),
		counter:   new(uint64),
	}

	for i := 0; i < bk; i++ {
		s.generator[uint64(i)] = &generator{sequence: new(uint64)}
	}
	return s, nil
}
