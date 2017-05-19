package hash

func HashStr(str string) int {
	var sum int
	for _, r := range str {
		sum += int(r) - 48
	}
	return sum % 1024
}
