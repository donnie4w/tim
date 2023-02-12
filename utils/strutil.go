/**
 * donnie4w@gmail.com  tim server
 */
package utils

import (
	"regexp"
	"strconv"
	"strings"
)

func Substr(str string, start, length int) string {
	if start < 0 || length < 0 {
		return str
	}
	rs := []rune(str)
	end := start + length
	return string(rs[start:end])
}

func IsStringEmpty(str string) bool {
	return strings.TrimSpace(str) == ""
}

func Atoi64(str string) int64 {
	i, _ := strconv.Atoi(str)
	return int64(i)
}

func Atoi32(str string) int32 {
	i, _ := strconv.Atoi(str)
	return int32(i)
}

func Atoi(str string) int {
	i, _ := strconv.Atoi(str)
	return i
}

func Atoi16(str string) int16 {
	i, _ := strconv.Atoi(str)
	return int16(i)
}

func IsNum(ss ...string) bool {
	if ss == nil || len(ss) == 0 {
		return false
	}
	pattern := "^\\d{0,}$"
	for _, s := range ss {
		b, err := regexp.MatchString(pattern, s)
		if err != nil || !b {
			return false
			break
		}
	}
	return true
}

func Join(ss ...string) string {
	if ss == nil || len(ss) == 0 {
		return ""
	}
	return strings.Join(ss, ",")
}
