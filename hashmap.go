package hashmap

import "errors"

type HashTableMode int

const (
	LoadFactorPricision float32       = 0.75
	SeperateChaining    HashTableMode = 0
	OpenAddressing      HashTableMode = 1
)

type HashFunc func(size int, key interface{}) (uint, uint)

type HashMap interface {
	Set(key interface{}, value interface{})
	Get(key interface{}) (interface{}, bool)
	Delete(key interface{}) error
}

func New(size int, mode HashTableMode, fn ...HashFunc) (HashMap, error) {
	if size < 1 {
		return nil, errors.New("bucket size should be over 0")
	}

	var hashFunc HashFunc
	if len(fn) > 0 {
		hashFunc = fn[0]
	} else {
		hashFunc = defaultHashFunc
	}

	switch mode {
	case SeperateChaining:
		return newHashMapWithSC(size, hashFunc)
	case OpenAddressing:
		return newHashMapWithOA(size, hashFunc)
	}

	return nil, errors.New("invalid hashtable mode")
}
