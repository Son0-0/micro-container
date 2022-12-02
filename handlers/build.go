package handlers

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func Build(Args []string) {
	fmt.Printf("build %v Image\n", Args[2])

	data, err := ioutil.ReadFile("Gockerfile")
	if err != nil {
		panic(err)
	}

	imageDir := "./images/"

	syscall.Umask(0)
	if err := os.Mkdir(imageDir+Args[2], 0777); err != nil {
		if strings.Contains(err.Error(), "file exists") {
			panic("Image already exists")
		}
	}

	command := strings.Split(string(data), "\n")

	sh, err := os.Create(imageDir + Args[2] + "/init.sh")
	Handle(err)
	defer sh.Close()

	workDir := imageDir + Args[2] + "/"

	for _, s := range command {
		if len(s) != 0 {
			cmd := strings.Split(string(s), " ")
			switch cmd[0] {
			case "FROM":
				dockerCommand := "docker export container > " + imageDir + Args[2] + "/image.tar"
				cmd := exec.Command("bash", "-c", dockerCommand)
				Handle(cmd.Run())

				unzipCommand := "tar xf " + imageDir + Args[2] + "/image.tar -C " + imageDir + Args[2] + "/"
				cmd = exec.Command("bash", "-c", unzipCommand)
				Handle(cmd.Run())

				Handle(os.Remove(imageDir + Args[2] + "/image.tar"))
			case "WORKDIR":
				syscall.Umask(0)
				Handle(os.Mkdir(imageDir+Args[2]+"/"+cmd[1], 0777))
				workDir = workDir + "/" + cmd[1] + "/"
				fmt.Fprintf(sh, "cd "+cmd[1]+"\n")
			case "CMD":
				tempCommand := ""
				for idx, c := range cmd[1:] {
					tempCommand += string(c)
					if idx != len(cmd[1:])-1 {
						tempCommand += " "
					} else {
						tempCommand += "\n"
					}
				}

				fmt.Fprintf(sh, string(tempCommand))
			case "COPY":
				data, err := ioutil.ReadFile(cmd[1])
				Handle(err)

				file, err := os.Create(workDir + cmd[1])
				defer file.Close()
				Handle(err)

				for _, s := range data {
					fmt.Fprintf(file, string(s))
				}
			}
		}
	}

	Handle(os.Chdir(imageDir))
	zipCommand := "tar cvf " + Args[2] + ".tar " + Args[2] + "/*"
	cmd := exec.Command("bash", "-c", zipCommand)
	Handle(cmd.Run())

	Handle(os.Chdir("../"))
	os.RemoveAll(imageDir + Args[2])
}
