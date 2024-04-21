package utils

func EmptyIsNull(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
