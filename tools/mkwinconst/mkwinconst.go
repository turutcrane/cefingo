package main

import (
	"bufio"
	"fmt"
	"os"
	strcase "github.com/stoewer/go-strcase"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		c := scanner.Text()
		
		fmt.Printf("\tWin%s = C.%s\n", strcase.UpperCamelCase(c), c)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
