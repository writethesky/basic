package utility

func ValidateStringRange(str string, minLength, maxLength int) bool {
	if minLength > len(str) {
		return false
	}
	if maxLength < len(str) {
		return false
	}
	return true
}

// ValidateStringHasTypes
// Multiple types are allowed
//     e.g. GenerateTypeNumber | GenerateTypeUpperLetter
func ValidateStringHasTypes(str string, exceptTypes GenerateType) (isOk bool) {
	for _, exceptType := range []GenerateType{GenerateTypeNumber, GenerateTypeUpperLetter, GenerateTypeLowerLetter, GenerateTypeSpecialCharacter} {
		exceptType = exceptTypes & exceptType
		if 0 != exceptType {
			isOk = validateStringHasType(str, exceptType)
			if !isOk {
				return
			}
		}
	}
	return
}

func validateStringHasType(str string, exceptType GenerateType) bool {
	var bytes []byte

	switch exceptType {
	case GenerateTypeNumber:
		bytes = numberLetters
	case GenerateTypeUpperLetter:
		bytes = upperLetters
	case GenerateTypeLowerLetter:
		bytes = lowerLetters
	case GenerateTypeSpecialCharacter:
		bytes = specialCharacters
	default:
		return false
	}

	for _, b := range bytes {
		for _, strB := range []byte(str) {
			if b == strB {
				return true
			}
		}
	}

	return false
}
