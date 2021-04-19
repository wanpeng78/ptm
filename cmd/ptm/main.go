package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"gitee.com/Stitchtor/ptm/filter"
	"gitee.com/Stitchtor/ptm/pkgm"
	"gitee.com/Stitchtor/ptm/syscheck"
	colorful "github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/urfave/cli/v2"
)

func main() {

	cfg := &Config{}
	app := &cli.App{
		Name:      "ptm",
		Usage:     "Packages tool mirror - 包管理镜像站设置工具",
		Copyright: "Copyright (c) " + strconv.Itoa(time.Now().Year()) + " Stitch-Zhang",
		Version:   fmt.Sprintf("%.2f", version),
		Flags:     cliParse(cfg),
		Action:    run(cfg),
	}
	app.UseShortOptionHandling = true
	app.HideHelpCommand = true
	sort.Sort(cli.FlagsByName(app.Flags))
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(cfg *Config) cli.ActionFunc {
	return func(c *cli.Context) error {
		ms := checkAndFetch(cfg.APIURL)
		osinfo, err := syscheck.SystemInfo()
		if err != nil {
			return err
		}
		osinfo.ShowInfo()
		urls := []string{}
		names := []string{}
		for _, v := range ms.Mirrors {
			urls = append(urls, v.URL)
			names = append(names, v.Name)
		}
		fmt.Println("正在测试延迟...")
		mirrorsInfo, err := filter.LatencyTest(urls, 5, true, 5, 3*time.Second, names...)
		if err != nil {
			colorful.Red("延迟测试失败")
			return nil
		}
		// choice mirrors
		mirrorsInfo.SortByDuration()
		mirrorsIdx := showMirrorsInfo(mirrorsInfo, *cfg)

		// translate mirror index to mirror
		targetMirrors := make(filter.Results, 0)
		for _, idx := range mirrorsIdx {
			targetMirrors = append(targetMirrors, mirrorsInfo[idx])
		}

		// prepare to write
		pm, err := pkgm.NewPKGM()
		if err != nil {
			colorful.Red("实例化包管理失败：%s\n", err.Error())
			os.Exit(0)
		}
		err = pm.WriteFile(targetMirrors)
		if err != nil {
			colorful.Red("写入镜像文件失败：%s\n", err.Error())
			os.Exit(0)
		}
		colorful.Green("写入成功，正在清理并更新镜像缓存")
		reader := pm.Refresh()
		displayRefresh(reader)
		colorful.Green("Everything is done!\n")
		return nil

	}
}

func showMirrorsInfo(mirrors filter.Results, cfg Config) (choices []int) {
	mirrorsTable := table.NewWriter()
	choicesTable := table.NewWriter()
	mirrorsTable.AppendHeader(table.Row{"#", "镜像名", "延迟(ms)", "地址", "备注"})
	choicesTable.AppendHeader(table.Row{"#", "镜像名", "延迟(ms)", "地址", "备注"})
	for idx, result := range mirrors {
		if result.Ping {
			mirrorsTable.AppendRow(table.Row{idx, result.Name, result.Latency.Milliseconds(), result.Address})
		} else {
			mirrorsTable.AppendRow(table.Row{idx, result.Name, result.Latency.Milliseconds(), result.Address, "远程服务器禁止ping，请以访问延迟为准"})
		}
	}
	fmt.Println(mirrorsTable.Render())

	// AutoConfig enable
	if cfg.AutoConfig {
		for idx, mirror := range mirrors {
			if mirror.Ping {
				choices = append(choices, idx)
			}
			if len(choices) >= int(cfg.MirrorCounts) {
				return
			}
		}
	}
	//Manual config mirrors
	choice := func() {
		var choicestr string
		appended := make(map[int]bool)
		fmt.Print("请输入所需要镜像站的序号(以,分割):")
		fmt.Scanf("%s", &choicestr)
		strs := strings.Split(choicestr, ",")
		for _, v := range strs {
			idx, err := strconv.Atoi(v)
			if err != nil || idx > len(mirrors) || idx < 0 {
				colorful.Red("序号错误:%d\n", idx)
				continue
			}
			if appended[idx] {
				colorful.Red("序号重复:%d", idx)
				continue
			}
			appended[idx] = true
			choices = append(choices, idx)
			if mirrors[idx].Ping {
				choicesTable.AppendRow(table.Row{idx, mirrors[idx].Name, mirrors[idx].Latency.Milliseconds(), mirrors[idx].Address})
			} else {
				choicesTable.AppendRow(table.Row{idx, mirrors[idx].Name, mirrors[idx].Latency.Milliseconds(), mirrors[idx].Address, "远程服务器禁止ping，请以访问延迟为准"})
			}
		}
	}
	for len(choices) == 0 {
		choice()
	}
	colorful.Green("已选的镜像:\n")
	fmt.Println(choicesTable.Render())
	return choices
}

func checkAndFetch(dataSourceURL string) *MirrorsOnline {
	// Check privileged
	if os.Geteuid() != 0 {
		colorful.Red("不支持非root用户使用,请以root权限运行")
		os.Exit(0)
	}
	// set ping privileged
	// Check update and fetch mirror source
	mirrorSource, err := fetchData(dataSourceURL)
	if err != nil {
		colorful.Red("无法从远程获取镜像数据:%s", err.Error())
		os.Exit(0)
	}
	return mirrorSource
}

func displayRefresh(buf *bufio.Reader) {
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			break
		} else {
			fmt.Print(line)
		}
	}
}
