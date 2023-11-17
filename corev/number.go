package corev

import (
	"fmt"
	"reflect"
	"strconv"
)

func NumberToInt(num interface{}) (int, error) {
	switch v := num.(type) {
	case int:
		return v, nil
	case int8:
		return int(v), nil
	case int16:
		return int(v), nil
	case int32:
		return int(v), nil
	case int64:
		return int(v), nil
	case uint:
		return int(v), nil
	case uint8:
		return int(v), nil
	case uint16:
		return int(v), nil
	case uint32:
		return int(v), nil
	case uint64:
		return int(v), nil
	case float32:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		return strconv.Atoi(v)
	case []byte:
		return strconv.Atoi(string(v))
	case bool:
		return IfInt(v, 1, 0), nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unsupported type: %s", reflect.TypeOf(num))
	}
}

func MustToInt(num interface{}) int {
	v, err := NumberToInt(num)
	if err != nil {
		panic(err)
	}
	return v
}

func AnyToString(v interface{}) string {
	switch v := v.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
