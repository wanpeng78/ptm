package apt

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"gitee.com/Stitchtor/ptm/filter"
	colorful "github.com/fatih/color"
)

// APT for apt in debian
type APT struct{}

// Version get apt's version
func (a *APT) Version() string {
	cmd := exec.Command("apt", "--version")
	r, _ := cmd.StdoutPipe()
	_ = cmd.Start()
	raw, _ := ioutil.ReadAll(r)
	return string(raw[4:])
}

// MirrorFilePath return mirror file absolute path
func (a APT) MirrorFilePath() string {
	return `/etc/apt/sources.list`
}

// BackupMirror back up local mirror file
// to source.list.bak
func (a APT) BackupMirror() error {
	source, err := ioutil.ReadFile(a.MirrorFilePath())
	if err != nil {
		return err
	}
	bakFileName := a.MirrorFilePath() + ".bak"
	err = ioutil.WriteFile(bakFileName, source, 0644)
	return err
}

// WriteFile write mirror to file
func (a *APT) WriteFile(mirrors filter.Results) error {
	notifaction :=
		`
	###
	### This file was generated by ptm
	### original file was renamed to sources.list.bak
	### author:Stitch-Zhang@gitee.com/stitchtor
	###
	### 默认注释了源码镜像以提高 apt update 速度，如有需要可自行取消注释
	`
	err := a.BackupMirror()
	general := "main restricted universe multiverse"
	if err != nil {
		return err
	}
	colorful.Green("mirror backup file : %s.bak", a.BackupMirror())
	// Open and delete file data
	file, err := os.OpenFile(a.MirrorFilePath(), os.O_RDWR|os.O_TRUNC, 0766)
	defer file.Close()
	if err != nil {
		return err
	}
	file.WriteString(notifaction + "\n")
	sysSymbol, err := sysSymbol()
	if err != nil {
		return err
	}
	for _, v := range mirrors {
		aptSymbols := aptSymbols(sysSymbol)
		for _, aptSymbol := range aptSymbols {
			mirror := fmt.Sprintf("%s %s %s", v.Address, aptSymbol, general)
			file.WriteString("deb " + mirror + "\n")
			file.WriteString("# deb-src " + mirror + "\n")
		}

	}
	return nil
}
