package main

import (
	"fmt"
)

func main() {
	done := make(chan bool, 1)

	N, a, c, x0 := 995300, 199061, 11, 15

	fmt.Printf("x(%7d) = %8d\n", 0, x0)
	outChan := GenerateRPN(N, a, c, x0, done)

	i := 1
	for n := range outChan {
		if i < 101 || n == x0 {
			fmt.Printf("x(%7d) = %8d\n", i, n)
		}

		if n == x0 {
			done <- true
			fmt.Printf("\nExpected period %d got %d\n", N, i)
			return
		}

		i += 1
	}
}