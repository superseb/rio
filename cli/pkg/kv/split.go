package kv

import "strings"

func Split(s, sep string) (string, string) {
	parts := strings.SplitN(s, sep, 2)
	return parts[0], safeIndex(parts, 1)
}

func safeIndex(parts []string, idx int) string {
	if len(parts) <= idx {
		return ""
	}
	return parts[idx]
}
