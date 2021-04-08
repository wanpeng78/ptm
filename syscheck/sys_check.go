package syscheck

import (
	"fmt"
	"log"

	"github.com/elastic/go-sysinfo"
	"github.com/elastic/go-sysinfo/types"
	colorful "github.com/fatih/color"
	"github.com/pkg/errors"
)

// SysType is a type of Linux distribution
type SysType string

// SysInfo is a warpper of types.OS
type SysInfo types.HostInfo

// ShowInfo display system info
func (s SysInfo) ShowInfo() {
	fmt.Println("System summary:")
	fmt.Print("System:")
	colorful.Green("%s@%s Patch:%d\n", s.OS.Name, s.OS.Version, s.OS.Patch)
	fmt.Print("Kernel:")
	colorful.Green("%s\n", s.KernelVersion)
	fmt.Print("Famliy:")
	colorful.Green("%s\n", s.OS.Family)
	fmt.Print("Archs:")
	colorful.Green("%s\n", s.Architecture)
}

const (
	// RedHat is a series of Linux which using yum for package manager
	// ex. Centos / RedHat Enterprise/ Fedora
	RedHat SysType = "redhat"

	// Debian is a series of Linux which using apt for package manager
	// ex. Ubuntu / Deepin / Kylin
	Debian SysType = "debian"

	// Arch is a series of Linux which using pacman for package manager
	// ex. Arch / Manjaro / KaOS
	Arch SysType = "arch"

	// UnSupported for others,which unsupported for now
	UnSupported SysType = "unsupported"
)

//Check to identify systemd type
func Check() string {
	info := sysinfo.Go()
	os := info.OS
	host, err := sysinfo.Host()
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Host:", host)
	fmt.Println("OS", os)
	return os
}

// SystemType return the type of the running linux
func SystemType() (SysType, error) {
	host, err := sysinfo.Host()
	if err != nil {
		return "", err
	}
	sysFamily := host.Info().OS.Family
	switch sysFamily {
	case string(Debian):
		return Debian, nil
	case string(RedHat):
		return RedHat, nil
	default:
		return UnSupported, errors.Errorf("unsupported type: %s", sysFamily)
	}
}

// SystemInfo return a full information of system
func SystemInfo() (SysInfo, error) {
	host, err := sysinfo.Host()
	if err != nil {
		return SysInfo{}, err
	}
	return SysInfo(host.Info()), nil
}
