package utility

import (
	"math/rand"
	"time"
)

type GenerateType uint

const (
	GenerateTypeNumber GenerateType = 1 << iota
	GenerateTypeUpperLetter
	GenerateTypeLowerLetter
	GenerateTypeSpecialCharacter
)

// GenerateRandomString
// You can set multiple generateTypes using "|"
//     e.g. GenerateTypeNumber | GenerateTypeUpperLetter
func GenerateRandomString(length int, generateTypes GenerateType) string {
	randomBytes := make([]byte, length)
	randomTypes := generateRandomTypes(length, generateTypes)
	for i, randomType := range randomTypes {
		randomBytes[i] = generateRandomByte(randomType)
	}
	return string(randomBytes)
}
func generateRandomTypes(length int, generateTypes GenerateType) (randomTypes []GenerateType) {
	randomTypes = make([]GenerateType, 0, length)
	for _, usedType := range []GenerateType{GenerateTypeNumber, GenerateTypeUpperLetter, GenerateTypeLowerLetter, GenerateTypeSpecialCharacter} {
		usedType = generateTypes & usedType
		if 0 != usedType {
			randomTypes = append(randomTypes, usedType)
		}
	}

	typeNum := len(randomTypes)
	if typeNum > length {
		randomTypes = randomTypes[0:length]
	}
	rand.Seed(time.Now().UnixNano())

	for i := typeNum; i < length; i++ {
		randomTypes = append(randomTypes, randomTypes[rand.Intn(typeNum)])
	}

	return
}

var specialCharacters = []byte("!@#$%^&*()/")
var numberLetters = []byte("0123456789")
var upperLetters = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
var lowerLetters = []byte("abcdefghijklmnopqrstuvwxyz")

func generateRandomByte(generateType GenerateType) byte {
	rand.Seed(time.Now().UnixNano())
	var bytes []byte
	switch generateType {
	case GenerateTypeNumber:
		bytes = numberLetters
	case GenerateTypeUpperLetter:
		bytes = upperLetters
	case GenerateTypeLowerLetter:
		bytes = lowerLetters
	case GenerateTypeSpecialCharacter:
		bytes = specialCharacters
	default:
		panic("Not an allowed generateType")
	}
	return bytes[rand.Intn(len(bytes))]
}
