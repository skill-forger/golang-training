package math_string

// ReverseString reverses a given string.
func ReverseString(s string) string {
	rns := []rune(s) // Convert string to rune slice to handle Unicode characters
	for i, j := 0, len(rns)-1; i < j; i, j = i+1, j-1 {
		rns[i], rns[j] = rns[j], rns[i]
	}
	return string(rns)
}

// IsPalindrome checks if a string is a palindrome.
func IsPalindrome(s string) bool {
	reversedS := ReverseString(s)
	return s == reversedS
}
