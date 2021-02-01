package main

import (
	"fmt"
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

	cmd := exec.Command("/usr/bin/oc", "apply", "-f", "/etc/origin/master/webconsole-config.yaml")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	services := []string{"api", "controllers"}
	for _, s := range services {
		fmt.Println(s)
		restartControlPlane(s)
	}

	time.Sleep(5)

	fmt.Println("Restarting web-console")
	deletePod("openshift-web-console")
	time.Sleep(5)
}
