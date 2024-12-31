package utils

import (
	"math"
	"unicode/utf8"
)

// Note that '_' char is formally invalid but is historically in use, especially on corpnets
func isValidDomainLabelChar(char rune) bool {
	if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '-' || char == '_' {
		return true
	}
	return false
}

func isLetterOrDigit(char rune) bool {
	if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') {
		return true
	}
	return false
}

func IsValidDomainName(name string) bool {
	end := utf8.RuneCountInString(name)
	if name == "" || end > math.MaxInt16 {
		return false
	}
	curPos := 0
	newPos := 0
	r := []rune(name)
	for {
		newPos = curPos
		for newPos < end {
			if r[newPos] == '.' {
				break
			}
			newPos++
		}
		if curPos == newPos || newPos-curPos > 63 || !isLetterOrDigit(rune(r[curPos])) {
			return false
		}
		curPos++
		for curPos < newPos {
			if !isValidDomainLabelChar(rune(r[curPos])) {
				return false
			}
			curPos++
		}
		curPos++
		if curPos >= end {
			break
		}
	}
	return true
}
