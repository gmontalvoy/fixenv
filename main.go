package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {

	if check("/etc/origin/master/master-config.yaml") {
		fmt.Println("Reconcile master-config.yaml")
		reconcileFiles("/etc/origin/master/master-config.yaml", 4)
	}

	if check("/etc/origin/master/webconsole-config.yaml") {
		fmt.Println("Reconcile webconsole-config.yaml")
		reconcileFiles("/etc/origin/master/webconsole-config.yaml", -1)
	}

	input, err := ioutil.ReadFile("/etc/origin/master/webconsole-config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	output := bytes.Replace(input, []byte("resourceVersion: \"12862\""), []byte("\n"), -1)
	if err = ioutil.WriteFile("/etc/origin/master/webconsole-config.yaml", output, 0666); err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("/usr/bin/oc", "apply", "-f", "/etc/origin/master/webconsole-config.yaml")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	apperr := cmd.Run()
	if apperr != nil {
		panic(err)
	}

	services := []string{"api", "controllers"}
	for _, s := range services {
		fmt.Println(s)
		restartControlPlane(s)
		time.Sleep(10 * time.Second)
	}

	fmt.Println("Restarting web-console")
	deletePod("openshift-web-console")
	time.Sleep(5)
}
