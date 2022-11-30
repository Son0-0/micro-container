package handlers

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func Run(Args []string) {
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
	}

	Handle(cmd.Run())
}

func Child(Args []string) {
	fmt.Printf("Runnig %v as PID %d\n", Args[2:], os.Getpid())

	cmd := exec.Command(Args[2], Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// set hostname
	Handle(syscall.Sethostname([]byte("container")))

	// chroot
	Handle(syscall.Chroot("/home/ubuntu/gocker"))
	Handle(syscall.Chdir("/"))

	// Mount proc
	Handle(syscall.Mount("proc", "proc", "proc", 0, ""))
	defer Handle(syscall.Unmount("proc", 0))

	Handle(cmd.Run())
}
