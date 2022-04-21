package strutil

func IsStrInList(str string, list []string) bool {
	for _, s := range list {
		if s == str {
			return true
		}
	}
	return false
}
