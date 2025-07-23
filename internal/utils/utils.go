package utils

func Reverse[S ~[]E, E any](s S) S {
	result := make(S, len(s))
	for i := 0; i < len(s); i++ {
		result[i] = s[len(s)-1-i]
	}

	return result
}
