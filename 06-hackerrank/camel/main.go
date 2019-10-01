package main

import (
	"fmt"
	"unicode"
)

func main() {
	var input string
	fmt.Scanf("%s", &input)
	fmt.Println(camelcase(input))
}

// Complete the camelcase function below.
func camelcase(s string) int32 {
	answer := 1
	for _, letter := range s {
		if unicode.IsUpper(letter) {
			answer++
		}
	}
	return int32(answer)
}
