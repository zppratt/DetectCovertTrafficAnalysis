package main

import "fmt"
import "flag"

func main() {
	var toPrint = flag.String("toPrint", "hello", "the string to print")
	fmt.Println(*toPrint)
}
