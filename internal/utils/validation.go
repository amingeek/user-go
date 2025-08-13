package utils

import "regexp"

func ValidatePhoneNumber(phone string) bool {
	// E.164 format validation
	pattern := `^\+[1-9]\d{1,14}$`
	matched, _ := regexp.MatchString(pattern, phone)
	return matched
}
