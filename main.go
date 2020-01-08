package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	fileName      = "changelist"
	conflicts_cmd = "git diff --name-only --diff-filter=U"
)

var ignoreList = []string{}

func main() {
	conflicts := runCmd(conflicts_cmd)
	if len(conflicts) == 0 {
		fmt.Println("No conflict")
		os.Exit(0)
	}

	loadChangeList()

	for _, v := range conflicts {
		if v == "" {
			continue
		}
		add := false
		var res []string
		opt := resolveOption(v)
		switch opt {
		case "ours":
			fmt.Println("ours: " + v)
			res = runCmd("git checkout --ours " + v)
			add = true
		case "theirs":
			fmt.Println("theirs: " + v)
			res = runCmd("git checkout --theirs " + v)
			add = true
		case "ignore":
			fmt.Println("ignore: " + v)
		case "warn":
		}
		if len(res) != 0 {
			fmt.Sprintln(strings.Join(res, "||"))
		} else {
			if add {
				runCmd("git add " + v)
			}
		}
	}
}

func loadChangeList() {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	br := bytes.NewReader(content)
	reader := bufio.NewReader(br)
	for {
		line, _, err := reader.ReadLine()
		if err != nil || io.EOF == err {
			break
		}
		if len(line) != 0 {
			ignoreList = append(ignoreList, string(line))
		}
	}
}

func runCmd(command string) []string {
	cmd := exec.Command("/bin/bash", "-c", command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	reader1 := bufio.NewReader(stderr)
	errstr := ""
	for {
		line, err2 := reader1.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		errstr += line + "\n"
	}
	if errstr != "" {
		fmt.Println(errstr)
		os.Exit(1)
	}

	result := make([]string, 0)
	reader2 := bufio.NewReader(stdout)
	for {
		line, err2 := reader2.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		result = append(result, line)
	}
	cmd.Wait()
	return result
}

func resolveOption(name string) string {
	ignore := false

	for _, v := range ignoreList {
		if strings.TrimSpace(v) == strings.TrimSpace(name) {
			ignore = true
			break
		}
	}
	if ignore {
		return "ignore"
	} else {
		return "theirs"
	}
}
