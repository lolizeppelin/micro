package utils

import (
	"fmt"
	"golang.org/x/exp/maps"
	"reflect"
	"sort"
	"strings"
)

// SortAbleSlice 可排序列表,从小到大
type SortAbleSlice[T SortAbleType] []T

func (x SortAbleSlice[T]) Len() int           { return len(x) }
func (x SortAbleSlice[T]) Less(i, j int) bool { return x[i] < x[j] }
func (x SortAbleSlice[T]) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

// IncludeInSlice 判断包含
func IncludeInSlice[T comparable](s []T, target T) bool {
	for _, v := range s {
		if v == target {
			return true
		}
	}
	return false

}

// InsertSlice 插入到最开始
func InsertSlice[T any](s []T, value T) []T {
	l := make([]T, len(s)+1)
	l[0] = value
	copy(l[1:], s)
	return l
}

// IntSliceToString int类型转换string
func IntSliceToString[T IntType](s []T) []string {
	output := make([]string, 0)
	for _, v := range s {
		output = append(output, IntToString(v))
	}
	return output
}

// MergeSlice 列表合并
func MergeSlice[T any](slices ...[]T) []T {
	var merged []T
	for _, slice := range slices {
		merged = append(merged, slice...)
	}
	return merged
}

// StringSliceToInt string 转换未int
func StringSliceToInt[T IntType](s []string) ([]T, error) {
	output := make([]T, 0)
	for _, v := range s {
		n, err := StringToInt(v)
		if err != nil {
			return nil, err
		}
		output = append(output, T(n))
	}
	return output, nil
}

// SortSlice 排序
func SortSlice[T SortAbleType](src []T, reverse ...bool) []T {
	s := CopySlice(src)
	if len(reverse) > 0 && reverse[0] {
		sort.Slice(s, func(i, j int) bool {
			return s[i] > s[j]
		})
	} else {
		sort.Slice(s, func(i, j int) bool {
			return s[i] < s[j]
		})
	}
	return s
}

// RemoveElementSlice 从列表中剔除元素(不创建新列表)
func RemoveElementSlice[T comparable](sp *[]T, target T) {
	s := *sp
	for i, item := range s {
		if item == target {
			s[i] = s[len(s)-1]
			*sp = s[:len(s)-1]
		}
	}
}

// MapValueToSlice map 值转换成列表
func MapValueToSlice[T1 comparable, T2 any](m map[T1]T2) []T2 {
	slice := make([]T2, 0)
	for _, v := range m {
		slice = append(slice, v)
	}
	return slice
}

// MapKeyToSlice map key转换成列表
func MapKeyToSlice[T1 comparable, T2 any](m map[T1]T2) []T1 {
	slice := make([]T1, 0)
	for k, _ := range m {
		slice = append(slice, k)
	}
	return slice
}

func Supported(supports []string, target string) bool {
	for _, value := range supports {
		if value == "*" || value == target {
			return true
		}
	}
	return false
}

// SliceReverse 列表倒转
func SliceReverse[T any](s []T) []T {
	newS := make([]T, len(s))
	for i, j := 0, len(s)-1; i <= j; i, j = i+1, j-1 {
		newS[i], newS[j] = s[j], s[i]
	}
	return newS
}

// SliceToMap 列表传字典,key通过传入的keyFunc 获取
func SliceToMap[T any, K comparable](s []T, keyFunc func(T) K) map[K]T {
	result := make(map[K]T)
	for _, item := range s {
		key := keyFunc(item)
		result[key] = item
	}
	return result
}

// HasIntersection  判断是否有交集
func HasIntersection[T comparable](slice1, slice2 []T) bool {
	set := make(map[T]struct{}) // struct{}不占空间
	for _, item := range slice1 {
		set[item] = struct{}{}
	}
	for _, item := range slice2 {
		if _, exists := set[item]; exists {
			return true
		}
	}
	return false
}

// SliceToMapByField 列表转字典，以列表中结构体的指定字段为key
func SliceToMapByField[T any, K comparable](s []T, field ...string) (map[K]T, error) {
	result := make(map[K]T)

	key := "ID"
	if len(field) > 0 {
		key = field[0]
	}

	for _, item := range s {
		val := reflect.ValueOf(item)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		// 获取字段值
		fieldVal := val.FieldByName(key)
		if !fieldVal.IsValid() {
			return nil, fmt.Errorf("field %s not found", key)
		}

		// 确保字段可以比较
		if !fieldVal.Type().Comparable() {
			return nil, fmt.Errorf("field %s is not comparable", key)
		}

		// 转换字段值为键类型
		k, ok := fieldVal.Interface().(K)
		if !ok {
			return nil, fmt.Errorf("field %v cannot be converted to key type", k)
		}

		result[k] = item
	}

	return result, nil
}

// SliceToMapByStringField 列表转字典，以列表中结构体的指定字符串字段为key
func SliceToMapByStringField[T any](s []T, key ...string) (map[string]T, error) {
	return SliceToMapByField[T, string](s, key...)
}

// RepeatSlice 实现python  [1] * n的功能
func RepeatSlice[T any](slice []T, times int) []T {
	if times <= 0 {
		return []T{}
	}

	// Calculate the new length for the resulting slice
	result := make([]T, 0, len(slice)*times)

	// Append the slice to itself 'times' times
	for i := 0; i < times; i++ {
		result = append(result, slice...)
	}
	return result
}

// SliceJoin 列表合并
func SliceJoin[T SortAbleType](values []T, quote string, split ...string) string {
	delimiter := ","
	if len(split) > 0 {
		delimiter = split[0]
	}
	var result []string
	for _, value := range values {
		result = append(result, fmt.Sprintf("%s%v%s", quote, value, quote))
	}
	return strings.Join(result, delimiter)
}

// StrSliceJoin 数字合并
func StrSliceJoin(values []string) string {
	return SliceJoin(values, `"`)
}

// NumberSliceJoin 数字合并
func NumberSliceJoin[T NumberType](values []T) string {
	return SliceJoin(values, "")
}

func Sets[T comparable](slices []T) []T {
	temp := map[T]struct{}{}
	for _, item := range slices {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
		}
	}
	return maps.Keys(temp)
}

func SplitSlice[T any](slice []T, total int) [][]T {
	if len(slice) < total {
		return [][]T{
			CopySlice(slice),
		}
	}
	var result [][]T
	for i := 0; i < len(slice); i += total {
		end := i + total
		if end > len(slice) {
			end = len(slice)
		}
		result = append(result, slice[i:end])
	}
	return result
}
