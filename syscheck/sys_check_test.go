package syscheck

import (
	"fmt"
	"testing"
)

func TestSystemType(t *testing.T) {
	sysType, err := SystemType()
	if err != nil {
		fmt.Println(err)
	}
	t.Log(sysType)
}

func TestShowInfo(t *testing.T) {
	info, err := SystemInfo()
	if err != nil {
		fmt.Println(err)
	}
	info.ShowInfo()
}
