package main

import (
	"log"

	"gitee.com/Stitchtor/ptm/syscheck"
)

func main() {
	sysinfo, err := syscheck.SystemInfo()
	if err != nil {
		log.Println(err)
		return
	}
	sysinfo.ShowInfo()
}
