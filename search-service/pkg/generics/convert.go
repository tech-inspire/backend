package generics

func Convert[A, B any](s []A, converter func(A) B) []B {
	out := make([]B, len(s))
	for i := range s {
		out[i] = converter(s[i])
	}

	return out
}

func OrConvert[A, B any](s *[]A, converter func(A) B) []B {
	if s == nil {
		return nil
	}

	return Convert(*s, converter)
}

func OrDefault[T any](v *T, defaultValue T) T {
	if v == nil {
		return defaultValue
	}
	return *v
}
