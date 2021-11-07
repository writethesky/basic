package utility

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateStringRange(t *testing.T) {
	min := 10
	max := 20
	for i := min; i <= max; i++ {
		isOk := ValidateStringRange(GenerateRandomString(i, GenerateTypeLowerLetter), min, max)
		assert.True(t, isOk)
	}

	for i := 0; i < min; i++ {
		isOk := ValidateStringRange(GenerateRandomString(i, GenerateTypeLowerLetter), min, max)
		assert.False(t, isOk)
	}

	for i := max + 1; i < max+100; i++ {
		isOk := ValidateStringRange(GenerateRandomString(i, GenerateTypeLowerLetter), min, max)
		assert.False(t, isOk)
	}

}

func TestValidateStringHasNumber(t *testing.T) {
	str := GenerateRandomString(10, GenerateTypeNumber|GenerateTypeUpperLetter)
	isOk := validateStringHasType(str, GenerateTypeNumber)
	assert.True(t, isOk)
}

func TestValidateStringHasUpperLetter(t *testing.T) {
	str := GenerateRandomString(10, GenerateTypeNumber|GenerateTypeUpperLetter)
	isOk := validateStringHasType(str, GenerateTypeUpperLetter)
	assert.True(t, isOk)

	str = GenerateRandomString(10, GenerateTypeNumber)
	isOk = validateStringHasType(str, GenerateTypeUpperLetter)
	assert.False(t, isOk)
}

func TestValidateStringHasLowerLetter(t *testing.T) {
	str := GenerateRandomString(10, GenerateTypeNumber|GenerateTypeLowerLetter)
	isOk := validateStringHasType(str, GenerateTypeLowerLetter)
	assert.True(t, isOk)

	str = GenerateRandomString(10, GenerateTypeNumber)
	isOk = validateStringHasType(str, GenerateTypeLowerLetter)
	assert.False(t, isOk)
}

func TestValidateStringHasSpecialCharacter(t *testing.T) {
	str := GenerateRandomString(10, GenerateTypeNumber|GenerateTypeSpecialCharacter)
	isOk := validateStringHasType(str, GenerateTypeSpecialCharacter)
	assert.True(t, isOk)

	str = GenerateRandomString(10, GenerateTypeNumber)
	isOk = validateStringHasType(str, GenerateTypeSpecialCharacter)
	assert.False(t, isOk)
}

func TestValidateStringHasUpperAndLowerLetter(t *testing.T) {
	str := GenerateRandomString(10, GenerateTypeNumber|GenerateTypeSpecialCharacter)
	isOk := ValidateStringHasTypes(str, GenerateTypeNumber|GenerateTypeSpecialCharacter)
	assert.True(t, isOk)

	str = GenerateRandomString(10, GenerateTypeNumber)
	isOk = ValidateStringHasTypes(str, GenerateTypeNumber|GenerateTypeSpecialCharacter)
	assert.False(t, isOk)
}
