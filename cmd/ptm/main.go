package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gitee.com/Stitchtor/ptm/syscheck"
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
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(cfg *Config) cli.ActionFunc {
	return func(c *cli.Context) error {
		osinfo, err := syscheck.SystemInfo()
		if err != nil {
			return err
		}
		osinfo.ShowInfo()
		fmt.Printf("%#v", *cfg)
		return nil

	}
}
