package zstrings

// CountByteRight counts the number of b from the
// right side of s until non b found.
// For example, CountByteRight("abcc",'c') returns 2.
func CountByteRight(s string, b byte) int {
	n := 0
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == b {
			n += 1
		} else {
			break
		}
	}
	return n
}

// CountByteLeft counts the number of b from the
// lest side of s until non b found.
// For example, CountByteLeft("aabc",'a') returns 2.
func CountByteLeft(s string, b byte) int {
	n := 0
	for i := 0; i <= len(s)-1; i++ {
		if s[i] == b {
			n += 1
		} else {
			break
		}
	}
	return n
}
