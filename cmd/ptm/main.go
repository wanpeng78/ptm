package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"gitee.com/Stitchtor/ptm/filter"
	"gitee.com/Stitchtor/ptm/syscheck"
	colorful "github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lorenzosaino/go-sysctl"
	"github.com/urfave/cli/v2"
)

func main() {
	// Check privileged
	if os.Geteuid() != 0 {
		colorful.Red("暂不支持非root用户使用,请以root权限运行")
		os.Exit(0)
	}
	// set ping privileged
	err := enablePing()
	if err != nil {
		colorful.Red("ping权限设置失败 %s", err.Error())
		os.Exit(0)
	}
	// Check update and fetch mirror source
	mirrorSource, err := fetchData()
	if err != nil {
		colorful.Red("无法从远程获取镜像数据:%s", err.Error())
	}
	if err != nil {
		colorful.Red("%s\n", err.Error())
	}
	cfg := &Config{}
	app := &cli.App{
		Name:      "ptm",
		Usage:     "Packages tool mirror - 包管理镜像站设置工具",
		Copyright: "Copyright (c) " + strconv.Itoa(time.Now().Year()) + " Stitch-Zhang",
		Version:   fmt.Sprintf("%.2f", version),
		Flags:     cliParse(cfg),
		Action:    run(cfg, mirrorSource),
	}
	app.UseShortOptionHandling = true
	app.HideHelpCommand = true
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(cfg *Config, ms *MirrorsOnline) cli.ActionFunc {
	return func(c *cli.Context) error {
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
		mirrorsInfo.SortByDuration()
		showMirrorsInfo(mirrorsInfo, false)
		return nil

	}
}

func showMirrorsInfo(mirrors filter.Results, onlyShow bool) []int {
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
	if onlyShow {
		return nil
	}
	var choicestr string
	var choices []int
	appended := make(map[int]bool)
	fmt.Print("请输入所需要镜像站的序号(以,分割):")
	fmt.Scanf("%s", &choicestr)
	strs := strings.Split(choicestr, ",")
	for _, v := range strs {
		idx, err := strconv.Atoi(v)
		if err != nil || idx > len(mirrors) {
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
	colorful.Green("已选的镜像:\n")
	fmt.Println(choicesTable.Render())
	return choices
}

// enable ipv4 ping
// More information @https://pkg.go.dev/github.com/go-ping/ping
func enablePing() error {
	return sysctl.Set("net.ipv4.ping_group_range", "0 2147483647")
}
