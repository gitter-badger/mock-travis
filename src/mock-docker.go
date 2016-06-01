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
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func addRepo(repoInfo string) {
	mockCfg := "/etc/mock/" + gyml("mock_travis.mock_config") + ".cfg"
	input, err := ioutil.ReadFile(mockCfg)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if line == `"""` {
			lines[i] = repoInfo
		}
	}

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(mockCfg, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func setTmpDir() {
	mkDir(tmpDir)
	if err := sh.Command("cp", "-r", shareDir, tmpDir+"/"+"SPEC").Run(); err != nil {
		boldColor("red", "Fail to setup sources to tmpdir.")
		os.Exit(1)
	}
	mkDir(tmpDir + "/" + "SRPM")
	mkDir(tmpDir + "/" + "RPM")
	mkDir(tmpDir + "/" + "debugInfo")
	mkDir(tmpDir + "/" + "source")
}

func setGit() {
	var (
		out []byte
		err error
	)
	gitUrl := "https://github.com/" + gyml("mock_travis.packages_buildrequires_git")
	boldColor("cyan", "Start setting git repository")
	_, _ = sh.Command("dnf", "-y", "install", "git").Output()
	out, err = sh.Command("git", "clone", gitUrl, tmpDir+"/"+"SPEC/GIT").Output()
	out = out[:0]
	if err != nil {
		boldColor("red", "Setting git repository failed.")
		os.Exit(1)
	}
	boldColor("green", "Setting git repository succeeded.")
}

func MockBuildSRPM(filePath string, f os.FileInfo, err error) error {
	var (
		outDown, outBuild []byte
		errDown, errBuild error
	)
	if filepath.Ext(f.Name()) == ".spec" {
		specDir := path.Dir(filePath)
		specFile := specDir + "/" + f.Name()
		boldColor("cyan", "Start downloading "+f.Name()+" source files")
		outDown, errDown = sh.Command("spectool", "-g", specFile, "-C", tmpDir+"/"+"source").Output()
		outDown = outDown[:0]
		if errDown != nil {
			boldColor("red", "Fail to download "+f.Name()+" source file")
			os.Exit(1)
		}
		boldColor("green", "Download "+f.Name()+" source succeeded.")
		boldColor("cyan", "Start building "+f.Name()+" SRPM")
		outBuild, errBuild = sh.Command("/usr/bin/mock",
			"-r",
			gyml("mock_travis.mock_config"),
			"--resultdir",
			tmpDir+"/"+"SRPM",
			"--buildsrpm",
			"--sources",
			tmpDir+"/"+"source",
			"--spec",
			specFile).Output()
		outBuild = outBuild[:0]
		if errBuild != nil {
			boldColor("red", "Build "+f.Name()+" SRPM failed")
		}
		boldColor("green", "Build "+f.Name()+" SRPM succeeded.")
	}
	return nil
}

func MockBuildRPM(filePath string, f os.FileInfo, err error) error {
	var (
		out []byte
	)
	if filepath.Ext(f.Name()) == ".rpm" {
		srpmDir := path.Dir(filePath)
		srpmFile := srpmDir + "/" + f.Name()
		boldColor("cyan", "Start building "+f.Name()+" binary RPM")
		out, err = sh.Command("/usr/bin/mock",
			"-r",
			gyml("mock_travis.mock_config"),
			"--resultdir",
			tmpDir+"/"+"RPM",
			"--rebuild",
			srpmFile,
		).Output()
		out = out[:0]
		if err != nil {
			boldColor("red", "Build "+f.Name()+" failed")
			rebuildList = append(rebuildList, srpmFile)
		} else {
			boldColor("green", "Build "+f.Name()+" succeeded.")
		}
	}
	return nil
}

func allBuild() {
	if gyml("mock_travis.packages_buildrequires_git") != "" {
		setGit()
	}
	specDir := tmpDir + "/" + "SPEC"
	srpmDir := tmpDir + "/" + "SRPM"
	_ = filepath.Walk(specDir, MockBuildSRPM)
	_ = filepath.Walk(srpmDir, MockBuildRPM)
	updateRepo()
	addRepo(localRepo)
	boldColor("yellow", "Start rebuilding for the binary RPMs built failed.")
	rebuildRPM()
	if len(stillFail) != 0 {
		boldColor("red", "Still build failed packages:")
		for i := 0; i < len(stillFail); i++ {
			boldColor("red", stillFail[i])
		}
		os.Exit(1)
	}
}

func rebuildRPM() {
	var (
		out []byte
		err error
	)
	for i := 0; i < cap(rebuildList); i++ {
		fileFullName := path.Base(rebuildList[i])
		boldColor("cyan", "Start rebuild "+fileFullName)
		out, err = sh.Command("/usr/bin/mock",
			"-r",
			gyml("mock_travis.mock_config"),
			"--resultdir",
			tmpDir+"/"+"RPM",
			"--rebuild",
			rebuildList[i],
		).Output()
		out = out[:0]
		if err != nil {
			boldColor("red", "Rebuild "+fileFullName+" failed")
			stillFail = append(stillFail, fileFullName)
		}
		boldColor("green", "Rebuild "+fileFullName+" succeeded.")
	}
}

func updateRepo() {
	var (
		out []byte
		err error
	)
	boldColor("cyan", "Start updating local repository")
	out, err = sh.Command("createrepo", tmpDir+"/"+"RPM").Output()
	out = out[:0]
	if err != nil {
		boldColor("red", "Update local repository failed")
		os.Exit(1)
	}
	boldColor("green", "Update local repository succeeded")
}

func initMock() {
	var (
		out []byte
		err error
	)
	boldColor("cyan", "Start setting up mock environment")
	out, err = sh.Command("/usr/bin/mock",
		"-r",
		gyml("mock_travis.mock_config"),
		"--init").Output()
	out = out[:0]
	if err != nil {
		boldColor("red", "Setup mock environment failed.")
	}
	boldColor("green", "Setup mock environment succeeded.")
}

func mockBuild() {
	setTmpDir()
	if gyml("mock_travis. packages_extra_repo") != "" {
		addRepo(extraRepo)
	}
	initMock()
	allBuild()
}
