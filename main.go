package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("invalid command")
	}
}

func run() {
	fmt.Printf("runnig %v\n", os.Args[2:])

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
	}

	handleError(cmd.Run())
}

func child() {
	fmt.Printf("child runnig %v as PID %d\n", os.Args[2:], os.Getpid())

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	handleError(syscall.Sethostname([]byte("container")))
	handleError(syscall.Chroot("/home/ubuntu/gocker"))
	handleError(syscall.Chdir("/"))
	handleError(syscall.Mount("proc", "proc", "proc", 0, ""))

	handleError(cmd.Run())

	handleError(syscall.Unmount("proc", 0))
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
