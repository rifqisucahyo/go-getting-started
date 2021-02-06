package controllers

import (
	"fmt"
	"math/rand"
	"strconv"
)

func InterfaceToFloat64(params interface{}) float64 {
	result, _ := strconv.ParseFloat(fmt.Sprint(params), 64)
	return result
}

func InterfaceToInt64(params interface{}) int64 {
	result, _ := strconv.ParseInt(fmt.Sprint(params), 10, 64)
	return result
}

func InArrayV2(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func substr(input string, start int, length int) string {
	asRunes := []rune(input)

	if start >= len(asRunes) {
		return ""
	}

	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}

	return string(asRunes[start : start+length])
}

func generateCode(n int) string {
	// const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const letterBytes = "0123456789"
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	)

	b := make([]byte, n)
	for i := 0; i < n; {
		if idx := int(rand.Int63() & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i++
		}
	}
	return string(b)
}
