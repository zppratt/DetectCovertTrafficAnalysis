package main

import (
	"io/ioutil"
	"os/exec"
)

func main() {
	if cmd, e := exec.Run("/bin/ls", nil, nil, exec.DevNull, exec.Pipe, exec.MergeWithStdout); e == nil {
		b, _ := ioutil.ReadAll(cmd.Stdout)
		println("output: " + string(b))
	}
}
