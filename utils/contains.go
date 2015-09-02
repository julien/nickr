package utils

func ContainsNickname(src []string, s string) bool {
	for _, i := range src {
		if i == s {
			return true
		}
	}
	return false
}

func AddNewNicknames(src []string, dst []string) []string {
	for _, i := range src {
		if !ContainsNickname(dst, i) {
			dst = append(dst, i)
		}
	}
	return dst
}
