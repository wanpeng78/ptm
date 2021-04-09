package apt

import (
	"io/ioutil"
	"os/exec"
)

type APT struct {
	version string
}

func (a *APT) Version() string {
	cmd := exec.Command("apt", "--version")
	r, _ := cmd.StdoutPipe()
	_ = cmd.Start()
	raw, _ := ioutil.ReadAll(r)
	a.version = string(raw[4:])
	return string(raw[4:])
}
