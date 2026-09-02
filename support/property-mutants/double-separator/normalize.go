package propertytests

func NormalizeKey(input string) string {
	result := make([]byte, 0, len(input))
	pendingSeparators := 0

	for i := 0; i < len(input); i++ {
		current := input[i]
		if !isASCIIAlphaNumeric(current) {
			if len(result) > 0 {
				pendingSeparators++
			}
			continue
		}
		for range pendingSeparators {
			result = append(result, '-')
		}
		pendingSeparators = 0
		result = append(result, toASCIILower(current))
	}
	return string(result)
}

func isASCIIAlphaNumeric(value byte) bool {
	return value >= 'a' && value <= 'z' ||
		value >= 'A' && value <= 'Z' ||
		value >= '0' && value <= '9'
}

func toASCIILower(value byte) byte {
	if value >= 'A' && value <= 'Z' {
		return value + ('a' - 'A')
	}
	return value
}
