package yum

import (
	"io/ioutil"
	"os/exec"
)

type YUM struct {
	version string
}

func (y *YUM) Version() string {
	cmd := exec.Command("yum", "--version")
	r, _ := cmd.StdoutPipe()
	_ = cmd.Start()
	raw, _ := ioutil.ReadAll(r)
	y.version = string(raw[:5])
	return string(raw[:5])
}
