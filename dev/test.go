package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

func main() {
	f, err := os.Open("input.dump.useful")
	bf := bufio.NewReader(f)
	fmt.Println(err)
	var line string
	//	var r = regexp.MustCompile("([0-9]{1,3}[\\.]){3}[0-9]{1,3}")
	var r = regexp.MustCompile("Source")
	//	var dest = regexp.MustCompile("Dest")
	//	var src = regexp.MustCompile("Source")
	for lnum := 0; lnum < 10; lnum++ {
		line, err = bf.ReadString(byte('\n'))
		fmt.Println(r.FindString(line))
	}
	fmt.Println(line)
}
