package pkgm

import (
	"bufio"
	"errors"

	"gitee.com/Stitchtor/ptm/filter"
	"gitee.com/Stitchtor/ptm/pkgm/apt"
	"gitee.com/Stitchtor/ptm/pkgm/yum"
	"gitee.com/Stitchtor/ptm/syscheck"
)

// PKGM for a package manager
// ex. yum or apt
type PKGM interface {
	Version() string
	MirrorFilePath() string
	WriteFile(filter.Results) error
	BackupMirror() error
	Refresh() *bufio.Reader
}

// NewPKGM return package manager
func NewPKGM() (PKGM, error) {
	sysType, err := syscheck.SystemType()
	if err != nil {
		return nil, err
	}
	switch sysType {
	case syscheck.Debian:
		return apt.New(), nil
	case syscheck.RedHat:
		return yum.New(), nil
	}
	return nil, errors.New("Unspported:" + string(sysType))
}

/*deb https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ hirsute main restricted universe multiverse
# deb-src https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ hirsute main restricted universe multiverse
deb https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ hirsute-updates main restricted universe multiverse
# deb-src https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ hirsute-updates main restricted universe multiverse
deb https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ hirsute-backports main restricted universe multiverse
# deb-src https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ hirsute-backports main restricted universe multiverse
deb https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ hirsute-security main restricted universe multiverse
# deb-src https://mirrors.tu
na.tsinghua.edu.cn/ubuntu/ hirsute-security main restricted universe multiverse*/
