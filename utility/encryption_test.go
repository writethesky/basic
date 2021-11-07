package utility

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptMd5WithNumber(t *testing.T) {
	for i := 10; i < 50; i++ {
		str := GenerateRandomString(i, GenerateTypeNumber)
		assert.Len(t, EncryptMd5(str), 32)
	}
}

func TestEncryptMd5WithLetter(t *testing.T) {
	for i := 10; i < 50; i++ {
		str := GenerateRandomString(i, GenerateTypeLowerLetter|GenerateTypeUpperLetter)
		assert.Len(t, EncryptMd5(str), 32)
	}
}

func TestEncryptMd5WithSpecialCharacter(t *testing.T) {
	for i := 10; i < 50; i++ {
		str := GenerateRandomString(i, GenerateTypeSpecialCharacter)
		assert.Len(t, EncryptMd5(str), 32)
	}
}
