package hashmap

import (
	"errors"
	"reflect"
	"sync"
)

type HashMapWithSC struct {
	mutex      sync.RWMutex
	fn         HashFunc
	elemCnt    int
	bucketSize int
	buckets    []*scBucket
	collCnt    int
}

type scBucket struct {
	hashKey uint
	key     interface{}
	value   interface{}
	next    *scBucket
}

func newHashMapWithSC(size int, fn HashFunc) (HashMap, error) {
	h := &HashMapWithSC{
		fn:         fn,
		bucketSize: getPrimeNum(size),
		buckets:    make([]*scBucket, getPrimeNum(size)),
	}
	return h, nil
}

func (h *HashMapWithSC) Set(key interface{}, value interface{}) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if h.getLoadFactor() > LoadFactorPricision {
		h.grow()
	}

	hKey, idx := h.fn(h.bucketSize, key)
	newBucket := &scBucket{hKey, key, value, nil}

	cursor := h.buckets[idx]
	if cursor == nil {
		h.buckets[idx] = newBucket
	} else {
		prev := h.buckets[idx]
		for cursor != nil {
			if reflect.DeepEqual(cursor.key, key) {
				cursor.value = value
				return
			}
			prev = cursor
			cursor = cursor.next
			h.collCnt++
		}
		prev.next = newBucket
	}
	h.elemCnt++
}

func (h *HashMapWithSC) GetCollisionCount() int {
	return h.collCnt
}

func (h *HashMapWithSC) Get(key interface{}) (interface{}, bool) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	_, idx := h.fn(h.bucketSize, key)
	cursor := h.buckets[idx]

	for cursor != nil {
		if reflect.DeepEqual(cursor.key, key) {
			return cursor.value, true
		}
		cursor = cursor.next
	}
	return nil, false
}

func (h *HashMapWithSC) Delete(key interface{}) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	_, idx := h.fn(h.bucketSize, key)
	cursor := h.buckets[idx]
	prev := h.buckets[idx]

	for cursor != nil {
		if reflect.DeepEqual(cursor.key, key) {
			if prev == cursor { // fisrt bucket
				h.buckets[idx] = cursor.next
				cursor = nil
			} else {
				prev.next = cursor.next
				cursor = nil
			}
			h.elemCnt--
			return nil
		}
		prev = cursor
		cursor = cursor.next
	}
	return errors.New("key now found")
}

func (h *HashMapWithSC) getLoadFactor() float32 {
	return float32(h.elemCnt) / float32(h.bucketSize)
}

func (h *HashMapWithSC) grow() {
	newBucketSize := getPrimeNum(h.bucketSize * 2)
	newBuckets := make([]*scBucket, newBucketSize)

	for i, bucket := range h.buckets {
		for bucket != nil {
			idx := bucket.hashKey % uint(newBucketSize)
			cursor := newBuckets[idx]
			newBucket := &scBucket{bucket.hashKey, bucket.key, bucket.value, nil}
			if cursor == nil {
				newBuckets[idx] = newBucket
			} else {
				prev := newBuckets[idx]
				for cursor != nil {
					prev = cursor
					cursor = cursor.next
				}
				prev.next = newBucket
			}
			bucket = bucket.next
		}
		h.buckets[i] = nil
	}

	h.bucketSize = newBucketSize
	h.buckets = newBuckets
}
