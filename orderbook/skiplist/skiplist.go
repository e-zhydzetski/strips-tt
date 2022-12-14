package skiplist

import (
	"math/rand"
	"sync"
)

func New[K any, V any](maxHeight int, compareFunc func(a K, b K) int) *SkipList[K, V] {
	return &SkipList[K, V]{
		maxRandHeightMask: 2 << (maxHeight - 1),
		compareFunc:       compareFunc,
		head:              new(Node[K, V]),
		nodePool: &sync.Pool{
			New: func() any {
				return new(Node[K, V])
			},
		},
	}
}

type SkipList[K any, V any] struct {
	maxRandHeightMask int
	compareFunc       func(a K, b K) int
	head              *Node[K, V]
	nodePool          *sync.Pool
}

type Node[K any, V any] struct {
	Key   K
	Value V
	Next  []*Node[K, V]
}

func (s *SkipList[K, V]) Set(key K, value V) {
	level := 0
	for r := rand.Intn(s.maxRandHeightMask); r&1 == 1; r >>= 1 {
		level++
	}
	if level >= len(s.head.Next) {
		level = len(s.head.Next)
		s.head.Next = append(s.head.Next, nil)
	}

	nn := s.nodePool.Get().(*Node[K, V])
	nn.Key = key
	nn.Value = value
	nn.Next = make([]*Node[K, V], level+1)

	cur := s.head
	for i := len(s.head.Next) - 1; i >= 0; i-- {
		for ; cur.Next[i] != nil; cur = cur.Next[i] {
			if s.compareFunc(cur.Next[i].Key, key) >= 0 {
				break
			}
		}
		if i <= level {
			nn.Next[i] = cur.Next[i]
			cur.Next[i] = nn
		}
	}
}
