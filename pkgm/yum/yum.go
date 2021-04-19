package yum

import (
	"bufio"
	"io/ioutil"
	"os"
	"os/exec"

	"gitee.com/Stitchtor/ptm/filter"
	colorful "github.com/fatih/color"
)

// YUM for CentOS's yum
type YUM struct{}

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
		file, err := os.OpenFile(filePath+repo.fileName, os.O_CREATE|os.O_RDWR, 0664)
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
func (y YUM) BackupMirror() error {
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

// Refresh cleanup cache and make new source cache
func (y YUM) Refresh() *bufio.Reader {
	cmd := exec.Command("yum", "makecache")
	out, _ := cmd.StdoutPipe()
	cmd.Start()
	return bufio.NewReader(out)
}
