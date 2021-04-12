package apt

import (
	"fmt"

	"gitee.com/Stitchtor/ptm/syscheck"
)

//for generate system char
func aptSymbols(sysmbol string) (result []string) {
	result = []string{
		sysmbol,
		sysmbol + "-updates",
		sysmbol + "-backports",
		sysmbol + "-security",
	}
	return
}

//short symbol for a version
//https://baike.baidu.com/item/ubuntu/155795?fr=aladdin#4_1
func sysSymbol() (string, error) {
	info, err := syscheck.SystemInfo()
	if err != nil {
		return "", err
	}
	major, minor := info.OS.Major, info.OS.Minor
	switch major {
	case 12:
		return "precise", nil
	case 14:
		return "trusty", nil
	case 16:
		return "xenial", nil
	case 18:
		return "bionic", nil
	case 20:
		if minor == 4 {
			return "focal", nil
		} else if minor == 10 {
			return "groovy", nil
		} else {
			return "", fmt.Errorf("unsupported version : %d.%d", major, minor)
		}
	case 21:
		return "hirsute", nil
	}
}
