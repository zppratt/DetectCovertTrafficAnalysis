package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// Call tshark and create the .useful file. Returns 0 on success, 1 on error.
func extract_useful() int {
	// Recieve custom destination for .useful
	var outputFileName []string

	if len(os.Args) > 1 {
		outputfilename = os.Args[1:]
	} else {
		outputfilename = []string{"."}
	}

	// Execute tshark and log output
	cmd := "tshark"
	args := []string{"-r", "input.dump", "-V", "foo.half.jpg"}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("Successfully ran tshark")


	return 0
}

// Extract the absolute time stamps for each IP pair, and write them to the corresponding file. Returns 0 on success, 1 on error.
func organize_ipd() int {
	return 0
}

// Display results (???). Returns 0 on succes, 1 on error. 
func extract_pid() int {
	return 0
}

func main() {
	extract_useful()
	organize_ipd()
	extract_pid()
}
