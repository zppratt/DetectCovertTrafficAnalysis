package main

import (
	"bufio"
	//	"bytes"
	"flag"
	"fmt"
	//	"io"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
//	"io/ioutil"
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
func extract_useful(path *string) int {
	fmt.Println("\nextract_useful():")
	fmt.Println("Reading " + *path)

	// Run tshark command, save .useful file
	fmt.Println("Creating " + *path + ".useful")
	err := exec.Command("sh", "-c", "/usr/bin/tshark -r "+*path+" -V | egrep \"Source:|Destination:|Time since reference or first frame:\" | grep -v \"Vmware\" > "+*path+".useful").Run()
	check(err)

	fmt.Println("Successfully created .useful file")
	return 0
}

// Extract the absolute time stamps for each IP pair, and write them to the corresponding . Returns 0 on success, 1 on error.
func organize_ipd(threads *int, path *string) int {
	fmt.Println("\norganize_ipd():")

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

	// Extract time stamps in .useful file using "threads" number of threads
	fmt.Println("Starting " + *path + ".useful parsing with " + strconv.Itoa(*threads) + " threads.")

	f, err := os.Open("input.dump.useful")
	check(err)
	defer f.Close()

	bf := bufio.NewReader(f)

	// Set up regex patterns
	var time_stamp = regexp.MustCompile("[0-9][.][0-9]{9}")
	var ip_addr = regexp.MustCompile("([0-9]{1,3}[\\.]){3}[0-9]{1,3}")
	var dest = regexp.MustCompile("Destination")
	var src = regexp.MustCompile("Source")

	var source_ip string
	var dest_ip string
	var time float64
	var ip_combo string

	var line string
	for lnum := 0; lnum < 99999; lnum++ {
		//		fmt.Println("Checking line ", lnum)
		line, err = bf.ReadString(byte('\n'))
		check(err)
		// If this line is a timestamp, set the timestamp
		if time_stamp.FindString(line) != "" {
			time, err = strconv.ParseFloat(time_stamp.FindString(line), 64)
			check(err)
		}
		// If this line contains the word "Source" and an IP address, set the source IP address
		if src.FindString(line) != "" && ip_addr.FindString(line) != "" {
			source_ip = ip_addr.FindString(line)
			//			fmt.Printf("Assigning %s as source\n", source_ip)
			/*
				If this line contains the word "Destination" and an IP address, set the destination IP address
				and create a unique mapping from source to destination
			*/
		} else if dest.FindString(line) != "" && ip_addr.FindString(line) != "" {
			dest_ip = ip_addr.FindString(line)
			//			fmt.Printf("Assigning %s as destination\n", dest_ip)
			if source_ip != "" {
				ipFiles[source_ip] = dest_ip
				check(err)
				//				fmt.Printf("Creating mapping %s_%s for timestamp %f\n", source_ip, dest_ip, time)
				ip_combo = source_ip + "_" + dest_ip
				timestamps[ip_combo] = append(timestamps[ip_combo], time)
				// Sort timestamps, inefficient, eh...
				sort.Float64s(timestamps[ip_combo])
			}
		}
	}

	for key, value := range ipFiles {
		fmt.Println("Key:", key, "Value:", value)
	}
	fmt.Println("***")
	var to_write string
	for key, value := range timestamps {
		for _, time := range timestamps[key] {
			to_write += strconv.FormatFloat(time, 'f', -1, 64) + "\n"
		}
		file, err := os.Create(key)
		check(err)
		fmt.Println("Key:", key, "Value:", value)
		_, err = file.WriteString(to_write)
		to_write=""
	}

	return 0
}

// Display results (???). Returns 0 on success, 1 on error.
func extract_pid() int {
	return 0
}

//func countLines(r io.Reader) (int, error) {
//	buf := make([]byte, 8196)
//	count := 0
//	lineSep := []byte{'\n'}
//
//	for {
//		c, err := r.Read(buf)
//		if err != nil && err != io.EOF {
//			return count, err
//		}
//
//		count += bytes.Count(buf[:c], lineSep)
//
//		if err == io.EOF {
//			break
//		}
//	}
//
//	return count, nil
//}

func main() {
	fmt.Println("\nmain():")

	// Custom name for the .useful file
	path, _ := os.Getwd()
	var dump = path + "/" + *flag.String("dump", "input.dump", "The name of the input file")
	fmt.Println(dump)
	// Number of threads to use as a command line argument
	var threads = flag.String("threads", "1", "The number of threads to use while parsing output")
	flag.Parse()

	// Convert threads flag value to int
	t, err := strconv.Atoi(*threads)
	check(err)

	fmt.Println(*threads + " thread(s) being used to parse .useful file.")

	// Run tshark on .dump file and create .useful file
	extract_useful(&dump)

	organize_ipd(&t, &dump)

	//extract_pid()
}
