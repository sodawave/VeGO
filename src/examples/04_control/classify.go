package main

import (
	"fmt"
	"os"
	"strconv"
)

func classify(n int) string {
	if n < 0 {
		return "neg"
	} else if n == 0 {
		return "zero"
	}
	return "pos"
}

func main() {
	n := 3
	for _, a := range os.Args[1:] {
		if a == "--" {
			continue
		}
		v, err := strconv.Atoi(a)
		if err == nil {
			n = v
			break
		}
	}
	acc := 0
	for i := 1; i <= n; i++ {
		acc = acc + i
	}
	fmt.Println(classify(n), acc)
}
