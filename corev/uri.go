package corev

import "strings"

/**
 * 获取url地址
 */
func GetQueryWithRune(uri string) string {
	return IfString(strings.ContainsRune(uri, '?'), IfString(strings.HasSuffix(uri, "?"), "", "&"), "?")
}

// 拼接路径
func GetPath(paths ...string) string {
	path := ""
	for _, p := range paths {
		if strings.HasSuffix(path, "/") && strings.HasPrefix(p, "/") {
			path += p[1:]
		} else if !strings.HasSuffix(path, "/") && !strings.HasPrefix(p, "/") {
			path += "/" + p
		} else {
			path += p
		}
	}
	return path
}
