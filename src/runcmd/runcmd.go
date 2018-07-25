package runcmd

import (
	"log"
	"os"
	"os/exec"
)

func Cmd(args ...string) {
	execmd := exec.Command(args[0], args[1:]...)
	execmd.Stderr = os.Stderr
	execmd.Stdout = os.Stdout
	execmd.Stdin = os.Stdin
	err := execmd.Run()
	if nil != err {
		log.Fatalln("Error running", execmd, ":", err)
	}
}
