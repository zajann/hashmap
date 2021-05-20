package hashmap

import (
	"testing"
)

const (
	startSize = 10000000
)

func BenchmarkSetWithMyHashMapSC(b *testing.B) {
	hm, err := New(startSize, SeperateChaining)
	if err != nil {
		b.Error(err)
	}

	for i := 0; i < b.N; i++ {
		hm.Set(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hm.Get(i)
	}
	b.Log(hm.(*HashMapWithSC).GetCollisionCount())
}

func BenchmarkSetWithMyHashMapOA(b *testing.B) {
	hm, err := New(startSize, OpenAddressing)
	if err != nil {
		b.Error(err)
	}

	for i := 0; i < b.N; i++ {
		hm.Set(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hm.Get(i)
	}
	b.Log(hm.(*HashMapWithOA).GetCollisionCount())
}

func BenchmarkSetWithBuiltInMap(b *testing.B) {
	m := make(map[int]int, startSize)

	for i := 0; i < b.N; i++ {
		m[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m[i]
	}
}
