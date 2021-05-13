package hashmap

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
)

func writeValue(buf *bytes.Buffer, val reflect.Value) {
	switch val.Kind() {
	case reflect.String:
		buf.WriteByte('"')
		buf.WriteString(val.String())
		buf.WriteByte('"')
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		buf.WriteString(strconv.FormatInt(val.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		buf.WriteString(strconv.FormatUint(val.Uint(), 10))
	case reflect.Float32, reflect.Float64:
		buf.WriteString(strconv.FormatFloat(val.Float(), 'E', -1, 64))
	case reflect.Bool:
		if val.Bool() {
			buf.WriteByte('t')
		} else {
			buf.WriteByte('f')
		}
	case reflect.Ptr:
		if !val.IsNil() || val.Type().Elem().Kind() == reflect.Struct {
			writeValue(buf, reflect.Indirect(val))
		} else {
			writeValue(buf, reflect.Zero(val.Type().Elem()))
		}
	case reflect.Array, reflect.Slice, reflect.Map, reflect.Struct, reflect.Interface:
		buf.WriteString(fmt.Sprintf("%#v", val))
	default:
		_, err := buf.WriteString(val.String())
		if err != nil {
			panic(fmt.Errorf("unsupported type %T", val))
		}
	}
}

func defaultHashFunc(size int, key interface{}) (uint, uint) {
	var buf bytes.Buffer
	writeValue(&buf, reflect.ValueOf(key))

	h := djb2Hash(&buf)

	return h, (h % uint(size))
}

func djb2Hash(buf *bytes.Buffer) uint {
	var h uint = 5381
	for _, r := range buf.Bytes() {
		h = (h << 5) + h + uint(r)
	}

	return h
}

func divideHash(size int, key uint) uint {
	return (key % uint(size+1)) + 1
}

func doubleHashFunc(size int, key interface{}, fn HashFunc, i int) uint {
	_, idx := fn(size, key)
	idx2 := divideHash(size, idx)

	return (idx + (uint(i) * idx2)) % uint(size)

}
