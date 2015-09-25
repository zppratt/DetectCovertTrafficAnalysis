/*

Zachary Paul Pratt
CS350
9/24/15

"By placing this statement in my work, I certify that I have read and understand the IPFW
Honor Code. I am fully aware of the following sections of the Honor Code: Extent of the
Honor Code, Responsibility of the Student and Penalty. This project or subject material
has not been used in another class by me or any other student. Finally, I certify that this
site is not for commercial purposes, which is a violation of the IPFW Responsible Use of
Computing (RUC) Policy."

ZACHARY P. PRATT, 9/23/15

*/

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
)

var ipFiles = map[string]string{}
var timestamps = map[string][]float64{}

func check(err error) {
	if err != nil {
		fmt.Println("Error occurred \"" + err.Error() + "\"")
		panic(err)
	}
}

// Call tshark and create the .useful file. Returns 0 on success, 1 on error.
func extract_useful(path *string) *string {
	useful := *path + ".useful"
	fmt.Println("Processing file \"" + useful + "\"")
	// Run tshark command, save .useful file
	fmt.Println("Creating \"" + useful + "\"")
	err := exec.Command("sh", "-c", "/usr/bin/tshark -r "+*path+" -V | egrep \"Source|Destination|Time since reference or first frame:\" | grep -v \"Vmware\" > "+useful).Run()
	check(err)
	fmt.Println("Successfully created \"" + useful + "\"")
	return &useful
}

// Extract the absolute time stamps for each IP pair, and write them to the corresponding . Returns 0 on success, 1 on error.
func organize_ipd(useful *string) {

	// Open .useful file
	//	file, err := os.Open(*path + ".useful")
	//	check(err)
	//	// Count lines in file
	//	lines, err := countLines(file)
	//	if err != nil {
	//		fmt.Println("Error counting line. Exiting.")
	//		os.Exit(1)
	//	}
	//
	//	fmt.Println("Number of lines in .useful file:", lines)
	//	linesPerProcessor := (lines / *threads)
	//	fmt.Println("Number of lines for each thread:", linesPerProcessor)

	f, err := os.Open(*useful)
	check(err)
	defer f.Close()

	bf := bufio.NewReader(f)

	// Set up regex patterns
	var time_stamp = regexp.MustCompile("[0-9][.][0-9]{9}")
	var ip_addr = regexp.MustCompile("([0-9]{1,3}[\\.]){3}[0-9]{1,3}")
	var src = regexp.MustCompile("Source")
	var dest = regexp.MustCompile("Destination")
	var port = regexp.MustCompile("[0-9]{5}")
	var src_port_regex = regexp.MustCompile("Source Port: [0-9]{5}")
	var dest_port_regex = regexp.MustCompile("Destination Port: [0-9]{5}")

	// Create string to hold relevant information
	var source_ip string
	var dest_ip string
	var source_port string
	var dest_port string
	var time float64
	var ip_combo_with_port string

	// Count lines in the file
	num_lines, err := countLines("input.dump.useful")
	check(err)
	fmt.Println("Processing " + strconv.Itoa(num_lines) + " lines of .useful file.")

	var line string
	for lnum := 0; lnum < num_lines; lnum++ {
		// Read the line
		line, err = bf.ReadString(byte('\n'))
		check(err)
		// If this line is a timestamp, set the timestamp
		if time_stamp.FindString(line) != "" {
			time, err = strconv.ParseFloat(time_stamp.FindString(line), 64)
			check(err)
		}
		if src_port_regex.FindString(line) != "" {
			source_port = port.FindString(line)
		}
		if dest_port_regex.FindString(line) != "" {
			dest_port = port.FindString(line)
		}
		// If this line contains the word "Source" and an IP address, set the source IP address
		if src.FindString(line) != "" && ip_addr.FindString(line) != "" {
			source_ip = ip_addr.FindString(line)
			//						fmt.Printf("Assigning %s/ as source\n", source_ip)
			/*
				If this line contains the word "Destination" and an IP address, set the destination IP address
				and create a unique mapping from source to destination
			*/
		} else if dest.FindString(line) != "" && ip_addr.FindString(line) != "" {
			dest_ip = ip_addr.FindString(line)
			if source_ip != "" && source_ip != "" && dest_port != "" {
				ipFiles[source_ip] = dest_ip
				check(err)
				ip_combo_with_port = source_ip + "_" + dest_ip + "_" + source_port + "_" + dest_port
				timestamps[ip_combo_with_port] = append(timestamps[ip_combo_with_port], time)
				// Sort timestamps, inefficient, eh...
				sort.Float64s(timestamps[ip_combo_with_port])
			}
		}
	}
	var to_write string
	var files_created int = 0
	for key, _ := range timestamps {
		for _, time := range timestamps[key] {
			to_write += strconv.FormatFloat(time, 'f', -1, 64) + "\n"
		}
		file, err := os.Create("output/" + key)
		files_created++
		check(err)
		_, err = file.WriteString(to_write)
		to_write = ""
	}
	fmt.Println(files_created, "files created for ip/port combination timestamps.")
}

// Display results (???). Returns 0 on success, 1 on error.
func extract_pid() int {
	return 0
}

func countLines(path string) (int, error) {
	r, err := os.Open(path)
	check(err)
	buf := make([]byte, 8196)
	count := 0
	lineSep := []byte{'\n'}
	for {
		c, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return count, err
		}
		count += bytes.Count(buf[:c], lineSep)
		if err == io.EOF {
			break
		}
	}
	return count, nil
}

func main() {
	// Custom name for the .useful file
	path, _ := os.Getwd()
	var dump = path + "/" + *flag.String("dump", "input.dump", "The name of the input file")
	flag.Parse()

	// Run tshark on .dump file and create .useful file
	useful := extract_useful(&dump)

	fmt.Println("Organizing \"" + *useful + "\"")
	organize_ipd(useful)

	//extract_pid()
}
