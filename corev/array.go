package corev

// SeachString

// ReverseStr 反正字符串
func ReverseStr(s string) string {
	if len(s) <= 1 {
		return s
	}
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}

// SplitStrCR 截取字符串
func SplitStrCR(s string, x rune, c int) string {
	for i, r := range s {
		if r == x {
			if c--; c == 0 {
				return s[:i]
			}
		}
	}
	return s
}

// ForEache ...
func ForEache(ds interface{}, next func(interface{}) (bool, error)) (err error) {
	ds2 := ds.(*[]interface{})
	for i := 0; i < len(*ds2); i++ {
		if n, e := next(&(*ds2)[i]); !n {
			break // 结束循环
		} else if e != nil {
			err = e
			break
		}
	}
	return
}

// FindArrayStr ...
func FindArrayStr(ds []string, s string) int {
	if len(ds) == 0 {
		return -1
	}
	for i, d := range ds {
		if d == s {
			return i
		}
	}
	return -1
}

// FindArrayString ...
func FindArrayString(ds []string, chk func(string) bool) (int, string) {
	if len(ds) == 0 {
		return -1, ""
	}
	for i, d := range ds {
		if chk(d) {
			return i, d
		}
	}
	return -1, ""
}

// FindArrayInterface ...
func FindArrayInterface(ds []interface{}, chk func(interface{}) bool) (int, interface{}) {
	if len(ds) == 0 {
		return -1, ""
	}
	for i, d := range ds {
		if chk(d) {
			return i, d
		}
	}
	return -1, ""
}
