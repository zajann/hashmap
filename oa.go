package hashmap

import (
	"errors"
	"reflect"
	"sync"
)

type HashMapWithOA struct {
	mutex      sync.RWMutex
	fn         HashFunc
	elemCnt    int
	bucketSize int
	buckets    []*oaBucket
	collCnt    int
}

var (
	delBucket oaBucket
)

type oaBucket struct {
	key   interface{}
	value interface{}
}

func newHashMapWithOA(size int, fn HashFunc) (HashMap, error) {
	h := &HashMapWithOA{
		fn:         fn,
		bucketSize: getPrimeNum(size),
		buckets:    make([]*oaBucket, getPrimeNum(size)),
	}
	return h, nil
}

func (h *HashMapWithOA) Set(key interface{}, value interface{}) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if h.getLoadFactor() > LoadFactorPricision {
		h.grow()
	}

	newBucket := &oaBucket{key, value}
	for i := 0; i < h.bucketSize; i++ {
		idx := doubleHashFunc(h.bucketSize, key, h.fn, i)
		if h.buckets[idx] == nil || h.buckets[idx] == &delBucket {
			h.buckets[idx] = newBucket
			break
		}
		if reflect.DeepEqual(h.buckets[idx].key, key) { // if same key, update
			h.buckets[idx].value = value
			return
		}
		h.collCnt++
	}
	h.elemCnt++
}

func (h *HashMapWithOA) GetCollisionCount() int {
	return h.collCnt
}

func (h *HashMapWithOA) Get(key interface{}) (interface{}, bool) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	for i := 0; i < h.bucketSize; i++ {
		idx := doubleHashFunc(h.bucketSize, key, h.fn, i)
		if h.buckets[idx] == nil {
			break
		}
		if reflect.DeepEqual(h.buckets[idx].key, key) {
			return h.buckets[idx].value, true
		}
	}
	return nil, false
}

func (h *HashMapWithOA) Delete(key interface{}) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for i := 0; i < h.bucketSize; i++ {
		idx := doubleHashFunc(h.bucketSize, key, h.fn, i)
		if h.buckets[idx] == nil {
			return errors.New("key not found")
		}
		if reflect.DeepEqual(h.buckets[idx].key, key) {
			h.buckets[idx] = &delBucket
			break
		}
	}
	h.elemCnt--
	return nil
}

func (h *HashMapWithOA) getLoadFactor() float32 {
	return float32(h.elemCnt) / float32(h.bucketSize)
}

func (h *HashMapWithOA) grow() {
	newBucketSize := getPrimeNum(h.bucketSize * 2)
	newBuckets := make([]*oaBucket, newBucketSize)

	for i, bucket := range h.buckets {
		if bucket == nil || bucket == &delBucket {
			continue
		}
		for i := 0; i < newBucketSize; i++ {
			idx := doubleHashFunc(newBucketSize, bucket.key, h.fn, i)
			if newBuckets[idx] == nil {
				newBuckets[idx] = &oaBucket{bucket.key, bucket.value}
				break
			}
		}
		h.buckets[i] = nil
	}
	h.bucketSize = newBucketSize
	h.buckets = newBuckets
}
