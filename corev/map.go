package corev

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func GetMapValue(m map[string]interface{}, key, sep string) (interface{}, bool) {
	keys := strings.Split(key, sep)
	var curr interface{} = m
	var next bool
	for _, kk := range keys {
		if vv, ok := curr.(map[string]interface{}); ok {
			curr, next = vv[kk]
		} else {
			next = false
		}
		if !next {
			return nil, false
		}
	}
	return curr, true
}

// AssertAsMap ...
func AssertAsMap(a interface{}) map[string]interface{} {
	r := reflect.ValueOf(a)
	if r.Kind().String() != "map" {
		panic(fmt.Sprintf("%v is not a map[string]interface{}", a))
	}

	res := make(map[string]interface{})
	tmp := r.MapKeys()
	for _, key := range tmp {
		res[key.String()] = r.MapIndex(key).Interface()
	}

	return res
}

// ConvertToStr ...
func ConvertToStr(a interface{}) string {
	res, err := json.Marshal(a)
	if err != nil {
		return ""
	}
	return string(res)
}

// ConvertToMap ...
func ConvertToMap(a interface{}) map[string]interface{} {
	t := reflect.TypeOf(a)
	v := reflect.ValueOf(a)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	if t.Kind() == reflect.Map {
		return AssertAsMap(a)
	} else if t.Kind() != reflect.Struct {
		panic("obj is not struct")
	}

	res := make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		key := field.Name
		tag := field.Tag.Get("json")
		if tag == "-" {
			continue
		}
		if idx := strings.Index(tag, ","); idx > 0 {
			tag = tag[:idx]
		}
		tag = strings.TrimSpace(tag)
		if tag == "" {
			tag = strings.ToLower(key)
		}
		res[tag] = v.Field(i).Interface()
	}
	return res
}

// Merge ...
func Merge(args ...interface{}) map[string]string {
	finalArg := make(map[string]string)
	for _, obj := range args {
		switch obj := obj.(type) {
		case map[string]*string:
			for key, value := range obj {
				if value != nil {
					finalArg[key] = *value
				}
			}
		default:
			byt, _ := json.Marshal(obj)
			arg := make(map[string]string)
			err := json.Unmarshal(byt, &arg)
			if err != nil {
				return finalArg
			}
			for key, value := range arg {
				if value != "" {
					finalArg[key] = value
				}
			}
		}
	}

	return finalArg
}

// MergeAll ... 非空内容不进行处理
func MergeAll(args ...interface{}) map[string]string {
	finalArg := make(map[string]string)
	for _, obj := range args {
		switch obj := obj.(type) {
		case map[string]*string:
			for key, value := range obj {
				if value != nil {
					finalArg[key] = *value
				}
			}
		default:
			byt, _ := json.Marshal(obj)
			arg := make(map[string]string)
			err := json.Unmarshal(byt, &arg)
			if err != nil {
				return finalArg
			}
			for key, value := range arg {
				finalArg[key] = value
			}
		}
	}

	return finalArg
}
