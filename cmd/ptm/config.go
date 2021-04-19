package main

import (
	"strconv"

	"github.com/urfave/cli/v2"
)

// Config for ptm
type Config struct {
	MirrorCounts    uint
	Interactive     bool
	OnlyShowMirrors bool
	AutoConfig      bool
	MirrorURL       string
	APIURL          string
}

const (
	mirrorFileURL = "https://gitee.com/Stitchtor/ptm/raw/master/raw/mirrors.json"
	version       = 0.1
	mirrorCounts  = 3
	link          = "https://gitee.com/Stitchtor/ptm"
)

func cliParse(c *Config) []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:        "interactive",
			Aliases:     []string{"it"},
			Usage:       "启用交互式配置",
			Destination: &c.Interactive,
			Value:       true,
			DefaultText: "true",
		},
		&cli.BoolFlag{
			Name:        "auto",
			Aliases:     []string{"at"},
			Usage:       "全自动配置",
			Destination: &c.AutoConfig,
			DefaultText: "false",
		},
		&cli.UintFlag{
			Name:        "mirrorCount",
			Aliases:     []string{"mc"},
			Usage:       "写入的镜像站数目",
			Value:       mirrorCounts,
			Destination: &c.MirrorCounts,
			DefaultText: strconv.Itoa(mirrorCounts),
		},
		&cli.StringFlag{
			Name:        "api",
			Usage:       "镜像数据文件地址",
			DefaultText: "https://gitee.com/Stitchtor/ptm/raw/master/raw/mirrors.json",
			Value:       mirrorFileURL,
			Destination: &c.APIURL,
		},
		&cli.StringFlag{
			Name:        "mirrorSites",
			Aliases:     []string{"mt"},
			Usage:       "自定义镜像地址，若启用则只会写入此一个镜像地址",
			Destination: &c.MirrorURL,
			DefaultText: "如http://mirror.neu.edu.cn",
		},
		&cli.BoolFlag{
			Name:        "onlyShowMirror",
			Aliases:     []string{"osm"},
			Usage:       "只显示镜像站点信的撒息,不进行操作",
			Destination: &c.OnlyShowMirrors,
		},
	}

}
