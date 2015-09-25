

package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

// A group contain the information of a pair
type PairEntry struct {
	time        float64
	source_ip   string
	dest_ip     string
	source_port string
	dest_port   string
}

// The list of entries
var pair_entries []PairEntry

const TIME = 0
const SOURCE_IP = 1
const DEST_IP = 2
const SOURCE_PORT = 3
const DEST_PORT = 4

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
	err := exec.Command("sh", "-c", "/usr/bin/tshark -r "+*path+" -V | egrep \"Source: ([0-9]{1,3}[\\.]){3}[0-9]{1,3}|Destination: ([0-9]{1,3}[\\.]){3}[0-9]{1,3}|Source Port: [0-9]{5}|Destination Port: [0-9]{5}|Time since\" > "+useful).Run()
	check(err)
	fmt.Println("Successfully created \"" + useful + "\"")
	return &useful
}

// Extract the absolute time stamps for each IP pair, and write them to the corresponding . Returns 0 on success, 1 on error.
func organize_ipd(useful *string) {

	// Set up regex patterns
	var time_stamp = regexp.MustCompile("[0-9][.][0-9]{9}")
	var src_ip_regex = regexp.MustCompile("Source:")
	var dest_ip_regex = regexp.MustCompile("Destination:")
	var src_port_regex = regexp.MustCompile("Source Port:")
	var dest_port_regex = regexp.MustCompile("Destination Port:")
	var ip = regexp.MustCompile("([0-9]{1,3}[\\.]){3}[0-9]{1,3}")
	var port = regexp.MustCompile("[0-9]{5}")

	var i int
	var entry int
	var found string

	// Open file, create reader
	file, err := os.Open(*useful)
	check(err)
	defer file.Close()
	r := bufio.NewReader(file)

	// Loop over lines in .useful file
	for {

//		fmt.Println("i=", i)

		l, _, err := r.ReadLine()

		if err != nil {
			break
		}

		line := string(l[:])

		// IP/Port info comes in sets of 5, loop over them and store information
		switch i % 5 {
		case TIME:
			// Initialize entry (necessary?)
			pair_entries[entry] = append(pair_entries[entry], PairEntry{0, "", "", "", ""})

			found = time_stamp.FindString(line)
			if found == "" {
				check(errors.New("Error reading timestamp on line" + strconv.Itoa(i)))
			} else {
				pair_entries[entry].time, err = strconv.ParseFloat(found, 64)
				check(err)
			}
		case SOURCE_IP:
			found = src_ip_regex.FindString(line)
			if found == "" {
				check(errors.New("Error reading src ip on line" + strconv.Itoa(i)))
			} else {
				pair_entries[entry].source_ip = ip.FindString(line)
			}
		case DEST_IP:
			found = dest_ip_regex.FindString(line)
			if found == "" {
				check(errors.New("Error reading dest ip on line" + strconv.Itoa(i)))
			} else {
				pair_entries[entry].dest_ip = ip.FindString(line)
			}
		case SOURCE_PORT:
			found = src_port_regex.FindString(line)
			if found == "" {
				check(errors.New("Error reading src port on line" + strconv.Itoa(i)))
			} else {
				pair_entries[entry].source_port = port.FindString(line)
			}
		case DEST_PORT:
			found = dest_port_regex.FindString(line)
			if found == "" {
				check(errors.New("Error reading dest port on line" + strconv.Itoa(i)))
			} else {
				pair_entries[entry].dest_port = port.FindString(line)
				// Iterate to next entry
				entry++
			}
		}

		i++

	}

	fmt.Printf("Done reading lines. %d lines found.\n", i)

	var count int

	for i = range pair_entries {
		count++
	}

	fmt.Printf("%d ip/port entries created.\n", count)

}

// Display results (???). Returns 0 on success, 1 on error.
func extract_pid() {
	var to_write string
	for key, _ := range timestamps {
		for _, time := range timestamps[key] {
			to_write += strconv.FormatFloat(time, 'f', -1, 64) + "\n"
		}
		file, err := os.Create("output/" + key + ".deltas")
		check(err)
		_, err = file.WriteString(to_write)
		to_write = ""
	}
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
