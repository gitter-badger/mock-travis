package main

import (
	"fmt"
	"github.com/codeskyblue/go-sh"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var (
	dockerName  string = "mock-build"
	dockerImage string = "nrechn/fedora-mock"
)

func boldColor(colorOption, msg string) {
	switch colorOption {
	case "red":
		_ = sh.Command("echo", "-e", `\n\e[31m\e[1m`+msg+`...\e[0m`).Run()

	case "green":
		_ = sh.Command("echo", "-e", `\n\e[32m\e[1m`+msg+`...\e[0m`).Run()

	case "yellow":
		_ = sh.Command("echo", "-e", `\n\e[33m\e[1m`+msg+`...\e[0m`).Run()

	case "cyan":
		_ = sh.Command("echo", "-e", `\n\e[36m\e[1m`+msg+`...\e[0m`).Run()

	default:
		fmt.Println("\n" + msg)
	}
}

func cleanDocker(name string) {
	var (
		out []byte
		err error
	)
	boldColor("cyan", "Start cleaning docker container...")
	out, err = sh.Command("docker", "rm", name).Output()
	out = out[:0]
	if err != nil {
		boldColor("red", "Clean docker container failed.")
		os.Exit(1)
	}
	boldColor("green", "Clean docker container succeeded.")
}

func currentLocation() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return dir
}

func gyml(arg string) string {
	viper.SetConfigName(".travis")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	return viper.GetString(arg)
}

func pullDocker() {
	var (
		out []byte
		err error
	)
	boldColor("cyan", "Start pulling "+dockerImage+" docker image...")
	out, err = sh.Command("docker", "pull", dockerImage).Output()
	out = out[:0]
	if err != nil {
		boldColor("red", "Pull "+dockerImage+" docker image failed.")
		os.Exit(1)
	}
	boldColor("green", "Pull "+dockerImage+" docker image succeeded.")
}

func runDocker(name string) {
	var (
		err            error
		volumeLocation string
	)
	volumeLocation = currentLocation() + "/:/home"
	if err = sh.Command("docker",
		"run",
		"--name",
		name,
		"--cap-add=SYS_ADMIN",
		"--privileged=true",
		"-v",
		volumeLocation,
		"-i",
		dockerImage,
		"/bin/fedora-mock").Run(); err != nil {
		boldColor("red",
			"OVERALL: Fail to build "+
				gyml("mock_travis.packages_name")+
				" and related build dependencies.")
		os.Exit(1)
	}
	boldColor("yellow",
		"OVERALL: Successfully build "+
			gyml("mock_travis.packages_name")+
			" and related build dependencies.")
}

func main() {
	pullDocker()
	runDocker(dockerName)
	cleanDocker(dockerName)
}
