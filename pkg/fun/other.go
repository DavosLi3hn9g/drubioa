package fun

//代替三元表达式
func If2String(condition bool, yes string, no string) string {
	if condition {
		return yes
	} else {
		return no
	}
}
func If2Int(condition bool, yes int, no int) int {
	if condition {
		return yes
	} else {
		return no
	}
}
func If2Bool(condition bool, yes bool, no bool) bool {
	if condition {
		return yes
	} else {
		return no
	}
}
