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

var sysVersion int

type repo struct {
	fileName string
	data     string
}

func init() {
	host, err := syscheck.SystemInfo()
	if err != nil {
		return
	}
	sysVersion = host.OS.Major
}

func repoGenerate(mirror filter.Result) (ret repo) {
	ret.fileName = fmt.Sprintf("CentOS-%s", mirror.Name)
	switch sysVersion {
	case 7:
		s := strings.ReplaceAll(centos7, "[mirror]", mirror.Address)
		ret.data = s
	case 8:
		s := strings.ReplaceAll(centos8, "[mirror]", mirror.Address)
		ret.data = s
	}
	return
}
