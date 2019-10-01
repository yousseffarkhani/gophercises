package main

import "fmt"

func main() {
	var length, rotate int32
	var s string
	fmt.Scanf("%d\n", &length)
	fmt.Scanf("%s\n", &s)
	fmt.Scanf("%d\n", &rotate)

	fmt.Println(length)
	fmt.Println(s)
	fmt.Println(rotate)
	fmt.Println(caesarCipher(s, rotate))
}

func caesarCipher(s string, k int32) string {
	var answer []rune
	k = k % 26
	for _, letter := range s {
		if letter >= 'a' && letter <= 'z' || letter >= 'A' && letter <= 'Z' {
			letter = (letter+k)%26 - 1
		}
		answer = append(answer, letter)
	}
	return string(answer)
}
