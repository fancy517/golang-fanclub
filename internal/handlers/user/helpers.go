package user

import (
	"math/rand"
	"strings"
)

const (
	LETTERS1 = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	LETTERS2 = "123456789"
	LETTERS3 = "abcdefghijklmnopqrstuvwxyz"
)

func generateActivationCode(codeLen int) string {
	res := make([]byte, codeLen)
	n := len(LETTERS1)
	for i := 0; i < codeLen; i++ {
		res[i] = LETTERS1[rand.Intn(n)]
	}
	return string(res)
}

func generatePasswordResetCode() string {
	res := make([]byte, 14)
	n := len(LETTERS3)
	for i := 0; i < 3; i++ {
		for j := 0; j < 4; j++ {
			res[j+5*i] = LETTERS3[rand.Intn(n)]
		}
	}
	res[4] = '-'
	res[9] = '-'
	return string(res)
}

func ValidateActivationCode(code string) bool {
	for i := 0; i < len(code); i++ {
		if strings.IndexByte(LETTERS1, code[i]) == -1 {
			return false
		}
	}
	return true
}
