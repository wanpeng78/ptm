package syscheck

import (
	"fmt"
	"log"

	"github.com/elastic/go-sysinfo"
)

// SysType is a type of Linux distribution
type SysType string
type A = string

const (
	// Fedora is a series of Linux which using yum for package manager
	// ex. Centos / RedHat Enterprise/ Fedora
	Fedora SysType = "Fedora"

	// Debian is a series of Linux which using yum for package manager
	// ex. Ubuntu / Deepin / Kylin

	Debian SysType = "Debian"
	// Arch is a series of Linux which using pacman for package manager
	// ex. Arch / Manjaro / KaOS

	Arch SysType = "Arch"
	// UnSupported for others,which unsupported for now
	UnSupported SysType = "UnSupported"

	a A = "DD"
)

//Check to identify systemd type
func Check() string {
	info := sysinfo.Go()
	os := info.OS
	host, err := sysinfo.Host()
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Host:", host.Info().OS.Version)
	fmt.Println("OS", os)
	return os
}
