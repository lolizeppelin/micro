package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
)

const (
	baseChars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	base      = 62
)

var (
	baseBytesChars = []byte(baseChars)
)

func IntToString[T IntType](value T) string {
	return strconv.Itoa(int(value))
}

func StringToInt(value string) (int, error) {
	return strconv.Atoi(value)
}

func StrToI32(value string) (int32, error) {
	v, err := StringToInt(value)
	if err != nil {
		return 0, err
	}
	if v > math.MaxInt {
		return 0, errors.New("value over max value of int32")
	}
	if v < math.MinInt {
		return 0, errors.New("value over min value of int32")
	}
	return int32(v), nil
}

func UnsafeStrToI32(value string) int32 {
	v, _ := StringToInt(value)
	return int32(v)
}

func JsonConvert(a, b any) error {
	buff, err := json.Marshal(a)
	if err != nil {
		return err
	}
	return json.Unmarshal(buff, b)
}

// ToString casts any type to a string type.
func ToString(i any) (string, error) {

	switch s := i.(type) {
	case float64:
		return strconv.FormatFloat(s, 'f', -1, 64), nil
	case float32:
		return strconv.FormatFloat(float64(s), 'f', -1, 32), nil
	case int:
		return strconv.Itoa(s), nil
	case int64:
		return strconv.FormatInt(s, 10), nil
	case int32:
		return strconv.FormatInt(int64(s), 10), nil
	case int16:
		return strconv.FormatInt(int64(s), 10), nil
	case int8:
		return strconv.FormatInt(int64(s), 10), nil
	case uint:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint64:
		return strconv.FormatUint(s, 10), nil
	case uint32:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(s), 10), nil
	case json.Number:
		return s.String(), nil
	case []byte:
		return string(s), nil
	case string:
		return s, nil
	case bool:
		return strconv.FormatBool(s), nil
	case nil:
		return "", nil
	case error:
		return s.Error(), nil
	default:
		v := reflect.ValueOf(s)
		switch v.Kind() {
		case reflect.String:
			return v.String(), nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return strconv.FormatInt(v.Int(), 10), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return strconv.FormatUint(v.Uint(), 0), nil
		case reflect.Float32, reflect.Float64:
			return strconv.FormatFloat(v.Float(), 'f', -1, 64), nil
		default:
			return "", fmt.Errorf("unable to cast %#v of type %T to string", i, i)
		}
	}
}

func UnsafeToString(i any) string {
	v, _ := ToString(i)
	return v
}
