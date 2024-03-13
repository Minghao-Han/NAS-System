package Utils

import (
	"strings"
)

var blockingWords = []string{"=", "where", "or", "and", "select", "delete", "update", "insert", "from", "values"}

func SQLInjectionDetector(str string) bool { //true=safe ,false = unsafe
	strLower := strings.ToLower(str)
	for _, blockingWord := range blockingWords {
		if strings.Contains(strLower, blockingWord) {
			return false
		}
	}
	return true
}
