package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	colorful "github.com/fatih/color"
)

// MirrorsOnline for fetching online mirror file to download the newest miroor
type MirrorsOnline struct {
	Maintainer string  `json:"maintainer"`
	LastUpdate string  `json:"last_update"`
	Version    float64 `json:"version"`
	Mirrors    []struct {
		Name      string   `json:"name"`
		Country   string   `json:"country"`
		URL       string   `json:"url"`
		Protocols []string `json:"protocols"`
	} `json:"mirrors"`
	ReleasedMsg string `json:"released_msg"`
}

func fetchData() (*MirrorsOnline, error) {
	colorful.Green("获取在线镜像数据.....")
	mo := new(MirrorsOnline)
	req, _ := http.NewRequest("GET", mirrorFileURL, nil)
	req.Header.Set("User-Agent", "ptm")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &mo)
	if err != nil {
		return nil, err
	}
	// Check update
	if mo.Version > version {
		colorful.Cyan("发现新版本: %.2f\n %s\n", mo.Version, link)
		colorful.Red("%s", mo.ReleasedMsg)
	}
	colorful.Green("数据获取成功")
	return mo, nil
}
