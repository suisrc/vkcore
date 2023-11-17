package corev

// IfString 判断
func IfString(condition bool, ifture, ifalse string) string {
	if condition {
		return ifture
	}
	return ifalse
}

// IfInt 判断
func IfInt(condition bool, ifture, ifalse int) int {
	if condition {
		return ifture
	}
	return ifalse
}

// IfBool 判断
func IfBool(condition bool, ifture, ifalse bool) bool {
	if condition {
		return ifture
	}
	return ifalse
}

// IfObject 判断
func IfObject(condition bool, ifture, ifalse interface{}) interface{} {
	if condition {
		return ifture
	}
	return ifalse
}

// IfFunc 判断
func IfFunc(condition bool, ifture, ifalse func() interface{}) interface{} {
	if condition && ifture != nil {
		return ifture()
	} else if !condition && ifalse != nil {
		return ifalse()
	}
	return nil
}

// IfExec 判断
func IfExec(condition bool, ifture, ifalse func()) {
	if condition && ifture != nil {
		ifture()
	} else if !condition && ifalse != nil {
		ifalse()
	}
}
