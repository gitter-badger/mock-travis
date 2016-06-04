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
	"github.com/codeskyblue/go-sh"
	"os"
)

func cleanDocker(name string) {
	var (
		err error
	)
	boldColor("cyan", "Start cleaning docker container...")
	_, err = sh.Command("docker", "rm", name).Output()
	if err != nil {
		boldColor("red", "Clean docker container failed.")
		os.Exit(1)
	}
	boldColor("green", "Clean docker container succeeded.")
}

func pullDocker() {
	var (
		err error
	)
	boldColor("cyan", "Start pulling "+dockerImage+" docker image...")
	_, err = sh.Command("docker", "pull", dockerImage).Output()
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
	volumeLocation = currentLocation() + "/:" + shareDir
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
		shareDir+"/mock-travis").Run(); err != nil {
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

func initDoc() {
	pullDocker()
	runDocker(dockerName)
	cleanDocker(dockerName)
}

func main() {
	if checkDocCon() {
		mockBuild()
	} else {
		initDoc()
	}

}
