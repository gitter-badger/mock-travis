// Copyright © 2016 nrechn <nrechn@gmail.com>
//
// This file is part of mock-travis.
//
// mock-travis is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// mock-travis is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with mock-travis. If not, see <http://www.gnu.org/licenses/>.
//

package main

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path"
	"path/filepath"
)

var (
	dockerName  string = "mock-build"
	dockerImage string = "nrechn/fedora-mock"
	shareDir    string = "/home"
	rebuildList []string
	stillFail   []string
	readTmpDir  string = "/var/tmp/birudo/" // Future compatibility
	tmpDir      string = path.Dir(readTmpDir + "/")
	extraRepo   string = `
[extra-local]
name=extra-local
baseurl=` + gyml("mock_travis.packages_extra_repo") + `
gpgcheck=0
"""
`

	localRepo string = `
[mock-local]
name=mock-local
baseurl=file://` + tmpDir + `/RPM/
gpgcheck=0
"""
`
)

func gyml(arg string) string {
	viper.SetConfigName(".travis")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/home/")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	return viper.GetString(arg)
}

func boldColor(colorOption, msg string) {
	switch colorOption {
	case "red":
		fmt.Println("\033[31m\033[1m" + msg + "\033[0m\033[39m")

	case "green":
		fmt.Println("\033[32m\033[1m" + msg + "\033[0m\033[39m")

	case "yellow":
		fmt.Println("\033[33m\033[1m" + msg + "\033[0m\033[39m")

	case "cyan":
		fmt.Println("\033[36m\033[1m" + msg + "\033[0m\033[39m")

	default:
		fmt.Println(msg)
	}
}

func checkDocCon() bool {
	if _, err := os.Stat("/.dockerinit"); err == nil {
		return true
	}
	return false
}

func currentLocation() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return dir
}

func mkDir(path string) {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		panic(err)
	}
}
