package filter

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/fatih/color"
)

func TestSpeedTest(t *testing.T) {
	var wg sync.WaitGroup
	urls := []string{
		"mirrors.tuna.tsinghua.edu.cn",
		"https://www.baidu.com",
		"www.google.com",
		"www.da.cc",
		"123.22.33.200",
		"tsing.mirrors.com/ubuntu/",
		"www",
	}
	c := make(chan Result, 1)
	go func() {
		for {
			select {
			case r := <-c:
				if r.Error != nil {
					color.Red("%s", r.Error)
					continue
				}
				color.Green("%s %d ms", r.Address, r.Latency.Milliseconds())
			}
		}
	}()
	errs := LatencyTestWithChan(urls, c, 5, 5*time.Second)
	if errs != nil {

	}
	wg.Wait()
	fmt.Println("---------")
}

func TestLatencyTestSlice(t *testing.T) {
	urls := []string{
		"mirrors.tuna.tsinghua.edu.cn",
		"https://www.baidu.com",
		"www.google.com",
		"www.da.cc",
		"123.22.33.200",
		"8.8.8.8",
		"tsing.mirrors.com/ubuntu/",
		"www",
		"https://mirrors.tuna.tsinghua.edu.cn/",
		"one.pantacor.com",
		"p.hjjjh.com",
	}
	results, err := LatencyTest(urls, 5, 5*time.Second)
	if err != nil {
		log.Panicln(err)
	}
	for _, v := range results {
		if v.Error != nil {
			color.Red("%s", v.Error)
			continue
		}
		color.Green("%s %d ms", v.Address, v.Latency.Milliseconds())
	}
}
