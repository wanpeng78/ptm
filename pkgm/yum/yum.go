package yum

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"gitee.com/Stitchtor/ptm/filter"
	colorful "github.com/fatih/color"
)

// YUM for CentOS's yum
type YUM struct {
}

// Version for get version of yum in localhost
func (y YUM) Version() string {
	cmd := exec.Command("yum", "--version")
	r, _ := cmd.StdoutPipe()
	_ = cmd.Start()
	raw, _ := ioutil.ReadAll(r)
	return string(raw[:6])
}

// WriteFile write files in mirror directory
// files depend on how many mirrors you choice
func (y YUM) WriteFile(mirrors filter.Results) error {
	filePath := y.MirrorFilePath()
	for _, mirror := range mirrors {
		repo := repoGenerate(mirror)
		file, err := os.OpenFile(filePath+repo.fileName, os.O_CREATE, 0664)
		if err != nil {
			return err
		}
		file.WriteString(repo.data)
		colorful.Green("Writed:%s\n", repo.fileName)
	}
	return nil
}

// BackupMirror for backing up local yum mirrors file
// Back`
func (y *YUM) BackupMirror() error {
	// we take advanced of yum sources
	// do not modify original yum source file
	// just add a new file which include of mirror source
	// new file name should be CentOS-[mirror name].repo
	return nil
}

//MirrorFilePath return yum source file path
func (y YUM) MirrorFilePath() string {
	return `/etc/yum.repos.d/`
}

func Hello() {
	mirrorFile := strings.Replace(centos7, "[Mirror]", "https://mirrors.tsinghua.tuna.com/CentOS/", -1)
	fmt.Println(mirrorFile)
}
