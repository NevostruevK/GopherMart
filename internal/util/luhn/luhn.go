package luhn

import (
	"regexp"
)

/*
	func Valid(s string) bool {
		var IsNumber = regexp.MustCompile(`^[0-9]+$`).MatchString
		if !IsNumber(s) {
			return false
		}
		sl := []byte(s)
		var luhn int
		pLast := len(sl)
		for i:=1; i<=len(sl);i++{
			cur := int(sl[pLast-i]) - 48
			if i%2 == 0 {
				cur = cur * 2
				if cur > 9{
					cur = cur%10 + cur/10
				}
			}
			luhn += cur
		}
		return luhn%10 == 0
	}
*/
func Valid(b []byte) bool {
	var IsNumber = regexp.MustCompile(`^[0-9]+$`).MatchString
	if !IsNumber(string(b)) {
		return false
	}
	var luhn int
	pLast := len(b)
	for i := 1; i <= len(b); i++ {
		cur := int(b[pLast-i]) - 48
		if i%2 == 0 {
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}
		luhn += cur
	}
	return luhn%10 == 0
}
