package yum

import (
	_ "embed"
	"fmt"
	"strings"

	"gitee.com/Stitchtor/ptm/filter"
	"gitee.com/Stitchtor/ptm/syscheck"
)

//embed feature only support Go 1.16+
//go:embed tmpl/centos7.tmpl
var centos7 string

//go:embed tmpl/centos8.tmpl
var centos8 string

type repo struct {
	fileName string
	data     string
}

// New return a YUM
func New() YUM {
	return YUM{}
}

func sysVersion() int {
	host, err := syscheck.SystemInfo()
	if err != nil {
		return 0
	}
	return host.OS.Major
}

func repoGenerate(mirror filter.Result) (ret repo) {
	ret.fileName = fmt.Sprintf("CentOS-%s.repo", mirror.Name)
	mirrorSite := fmt.Sprintf("%scentos", mirror.Address)
	switch sysVersion() {
	case 7:
		s := strings.ReplaceAll(centos7, "[Mirror]", mirrorSite)
		s = strings.ReplaceAll(s, "Name", mirror.Name)
		ret.data = s
	case 8:
		s := strings.ReplaceAll(centos8, "[Mirror]", mirrorSite)
		s = strings.ReplaceAll(s, "Name", mirror.Name)
		ret.data = s
	}
	return
}
