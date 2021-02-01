package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"time"
)

func reconcileFiles(files string, conc int) {

	input, err := ioutil.ReadFile(files)
	if err != nil {
		log.Fatal(err)
	}

	output := bytes.Replace(input, []byte("generic"), []byte("dynamic"), conc)
	if err = ioutil.WriteFile(files, output, 0666); err != nil {
		log.Fatal(err)
	}
}

func restartControlPlane(svc string) {

	cmd := exec.Command("/usr/local/bin/master-restart", svc)
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	time.Sleep(5)
}

func check(filecheck string) bool {

	file, err := os.Open(filecheck)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	r, err := regexp.Compile("dynamic")
	if err != nil {
		log.Fatal(err)
	}

	var a int = 0
	for scanner.Scan() {
		if r.MatchString(scanner.Text()) {
			a++
		}
	}

	if a > 4 {
		fmt.Printf("%v already reconciled, skip...\n", filecheck)
		return false
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return true
}

func deletePod(namespace string) {
	cmd := exec.Command("/usr/bin/oc", "delete", "pod", "--all", "-n", namespace)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
