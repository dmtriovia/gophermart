package validatef

import (
	"fmt"
	"regexp"
)

func IsMatchesTemplate(
	addr string,
	pattern string,
) (bool, error) {
	res, err := MatchString(pattern, addr)
	if err != nil {
		return false, err
	}

	return res, err
}

func MatchString(pattern string, s string) (bool, error) {
	re, err := regexp.Compile(pattern)
	if err == nil {
		return re.MatchString(s), nil
	}

	return false, fmt.Errorf("MatchString: %w", err)
}

const tmp = 10

func IsValidLuna(number int) bool {
	return (number%10+checksum(number/tmp))%tmp == 0
}

func checksum(number int) int {
	var luhn int

	for i := 0; number > 0; i++ {
		cur := number % tmp

		if i%2 == 0 {
			cur *= 2
			if cur > tmp-1 {
				cur = cur%tmp + cur/tmp
			}
		}

		luhn += cur
		number /= tmp
	}

	return luhn % tmp
}
