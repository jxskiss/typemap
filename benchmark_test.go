package phimap

import (
	"sync"
	"testing"
	"unsafe"
)

func Benchmark_Concurrent_StdMap_Get_NoLock(b *testing.B) {
	m := make(map[uintptr]uintptr)
	typPtrs := fillMap(func(k, v uintptr) { m[k] = v })

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for _, ptr := range typPtrs {
				_ = m[ptr]
			}
		}
	})
}

func Benchmark_Concurrent_StdMap_Get_RWMutex(b *testing.B) {
	var mu sync.RWMutex
	m := make(map[uintptr]uintptr)
	typPtrs := fillMap(func(k, v uintptr) { m[k] = v })

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for _, ptr := range typPtrs {
				mu.RLock()
				_ = m[ptr]
				mu.RUnlock()
			}
		}
	})
}

func Benchmark_Concurrent_SyncMap_Get(b *testing.B) {
	m := sync.Map{}
	typPtrs := fillMap(func(k, v uintptr) { m.Store(k, v) })

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for _, ptr := range typPtrs {
				got, _ := m.Load(ptr)
				_ = got.(uintptr)
			}
		}
	})
}

func Benchmark_Concurrent_Slice_Index(b *testing.B) {
	slice := make([]uintptr, 12)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < 12; i++ {
				_ = slice[i]
			}
		}
	})
}

func Benchmark_Concurrent_TypeMap_Get(b *testing.B) {
	m := NewTypeMap[uintptr]()
	typPtrs := fillMap(func(k, v uintptr) {
		_, _ = m.SetByUintptr(k, func() (uintptr, error) { return v, nil })
	})
	m.calibrate(true)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for _, ptr := range typPtrs {
				m.GetByUintptr(ptr)
			}
		}
	})
}

func fillMap(setfunc func(k, v uintptr)) []uintptr {
	var values = []any{
		TestType1{},
		TestType2{},
		TestType3{},
		TestType4{},
		TestType5{},
		TestType6{},
		TestType7{},
		TestType8{},
		TestType9{},
		TestType10{},
		TestType11{},
		TestType12{},
	}
	ptrs := make([]uintptr, 0, len(values))
	for _, val := range values {
		typPtr := (*(*[2]uintptr)(unsafe.Pointer(&val)))[1]
		setfunc(typPtr, typPtr)
		ptrs = append(ptrs, typPtr)
	}
	return ptrs
}
