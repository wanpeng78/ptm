package main

import (
	"strconv"

	"github.com/urfave/cli/v2"
)

// Config for ptm
type Config struct {
	MirrorsNums     int
	Interactive     bool
	OnlyShowMirrors bool
	MirrorURL       string
	APIURL          string
}

const (
	mirrorFileURL = "https://gitee.com/Stitchtor/ptm/raw/master/raw/mirrors.json"
	version       = 0.1
	mirrors       = 3
	link          = "https://gitee.com/Stitchtor/ptm"
)

func cliParse(c *Config) []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:        "interactive",
			Aliases:     []string{"it"},
			Usage:       "启用交互式配置",
			Destination: &c.Interactive,
			DefaultText: "否",
		},
		&cli.IntFlag{
			Name:        "mirrorcount",
			Aliases:     []string{"mc"},
			Usage:       "写入的镜像站数目",
			Destination: &c.MirrorsNums,
			DefaultText: strconv.Itoa(mirrors),
		},
		&cli.StringFlag{
			Name:        "api",
			Usage:       "镜像数据文件地",
			DefaultText: "https://gitee.com/Stitchtor/ptm/raw/master/raw/mirrors.json",
			Destination: &c.APIURL,
		},
		&cli.StringFlag{
			Name:        "mirrorsite",
			Aliases:     []string{"mt"},
			Usage:       "自定义镜像地址，若启用则只会写入此一个镜像地址",
			Destination: &c.MirrorURL,
			DefaultText: "如http://mirror.neu.edu.cn",
		},
		&cli.BoolFlag{
			Name:        "onlyShowMirror",
			Aliases:     []string{"osm"},
			Usage:       "只显示镜像站点信息,不进行操作",
			Destination: &c.OnlyShowMirrors,
			DefaultText: "否",
		},
	}

}
