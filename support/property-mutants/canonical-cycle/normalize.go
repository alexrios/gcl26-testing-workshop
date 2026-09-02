package propertytests

func NormalizeKey(input string) string {
	result := normalizeBase(input)
	if result == "a" {
		return "b"
	}
	if result == "b" {
		return "a"
	}
	return result
}

func normalizeBase(input string) string {
	result := make([]byte, 0, len(input))
	pendingSeparator := false

	for i := 0; i < len(input); i++ {
		current := input[i]
		if !isASCIIAlphaNumeric(current) {
			pendingSeparator = true
			continue
		}
		if pendingSeparator && len(result) > 0 {
			result = append(result, '-')
		}
		pendingSeparator = false
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
