package utility

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleGenerateRandomString() {

	GenerateRandomString(10, GenerateTypeNumber)
	// return => "1965016892"

	GenerateRandomString(10, GenerateTypeNumber|GenerateTypeUpperLetter)
	// return => "1I65A16PZ2"

	GenerateRandomString(10, GenerateTypeNumber|GenerateTypeUpperLetter|GenerateTypeLowerLetter)
	// return	=> "aI65A16lZ2"

	GenerateRandomString(10, GenerateTypeNumber|GenerateTypeUpperLetter|GenerateTypeLowerLetter|GenerateTypeSpecialCharacter)
	// return	=> "aI!5A1@lZ2"
}

func TestGenerateRandomString(t *testing.T) {
	randomLength := 11
	randomString := GenerateRandomString(randomLength, GenerateTypeNumber|GenerateTypeUpperLetter|GenerateTypeLowerLetter)
	assert.Len(t, randomString, randomLength)
}

func TestGenerateRandomTypesOnlyOne(t *testing.T) {
	for _, optionType := range []GenerateType{GenerateTypeNumber, GenerateTypeUpperLetter, GenerateTypeLowerLetter, GenerateTypeSpecialCharacter} {
		randomTypes := generateRandomTypes(10, optionType)
		assert.Equal(t, 10, len(randomTypes))
		for _, randomType := range randomTypes {
			assert.Equal(t, optionType, randomType)
		}
	}
}

func TestGenerateRandomTypesMultiType(t *testing.T) {
	randomTypes := generateRandomTypes(10, GenerateTypeNumber|GenerateTypeUpperLetter|GenerateTypeLowerLetter)
	assert.Equal(t, 10, len(randomTypes))
	groupType := make(map[GenerateType]bool)
	for _, randomType := range randomTypes {
		groupType[randomType] = true
	}
	assert.Equal(t, 3, len(groupType))
}

func TestGenerateRandomByteNumber(t *testing.T) {
	b := generateRandomByte(GenerateTypeNumber)
	// 0-9 : 48-57
	assert.GreaterOrEqual(t, b, byte(48))
	assert.LessOrEqual(t, b, byte(57))
}

func TestGenerateRandomByteUpperLetter(t *testing.T) {
	b := generateRandomByte(GenerateTypeUpperLetter)
	// A-Z: 65-90
	assert.GreaterOrEqual(t, b, byte(65))
	assert.LessOrEqual(t, b, byte(90))
}

func TestGenerateRandomByteLowerLetter(t *testing.T) {
	b := generateRandomByte(GenerateTypeLowerLetter)
	// a-z: 97-122
	assert.GreaterOrEqual(t, b, byte(97))
	assert.LessOrEqual(t, b, byte(122))
}

func TestGenerateRandomByteSpecialCharacter(t *testing.T) {
	b := generateRandomByte(GenerateTypeSpecialCharacter)
	assert.Contains(t, specialCharacters, b)
}
