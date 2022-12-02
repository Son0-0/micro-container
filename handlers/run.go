package handlers

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// --name <container_name> <image_name>
func Run(Args []string) {

	// check container name
	if Args[0] != "--name" {
		panic("Invalid container name...")
	}

	containerDir := "./containers/"
	imageDir := "./images/"

	// create container directory
	syscall.Umask(0)
	Handle(os.Mkdir(containerDir+Args[1], 0777))

	// check images and unzip images to containers directory
	imageUnzipCommand := "tar xf " + imageDir + Args[2] + ".tar -C" + containerDir + Args[1] + "/"
	cmd := exec.Command("bash", "-c", imageUnzipCommand)
	Handle(cmd.Run())

	cmd = exec.Command("/proc/self/exe", append([]string{"child"}, Args[1:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
	}

	Handle(cmd.Run())
}

func Child(Args []string) {
	// Args[2]: container name
	// Args[3]: image name
	fmt.Printf("Runnig %v as PID %d\n", Args[2:], os.Getpid())

	// init command
	cmd := exec.Command("/bin/bash", "-c", "sh init.sh")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// set hostname
	Handle(syscall.Sethostname([]byte(Args[2])))

	// chroot
	Handle(syscall.Chroot("./containers/" + Args[2] + "/" + Args[3]))
	Handle(syscall.Chdir("/"))

	// Mount proc
	Handle(syscall.Mount("proc", "proc", "proc", 0, ""))
	defer Handle(syscall.Unmount("proc", 0))

	Handle(cmd.Run())
}
