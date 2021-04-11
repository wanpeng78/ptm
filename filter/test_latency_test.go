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
		"https://mirrors.huaweicloud.com/",
		"https://mirrors.cloud.tencent.com/",
		"http://mirrors.163.com/",
		"http://mirrors.sohu.com/",
		"http://mirrors.yun-idc.com/",
		"https://mirrors.tuna.tsinghua.edu.cn/",
		"http://mirrors.ustc.edu.cn/",
		"http://mirror.neu.edu.cn/",
		"https://mirror.bjtu.edu.cn/",
		"http://mirrors.zju.edu.cn/",
		"http://mirrors.hust.edu.cn/",
		"https://mirrors.njupt.edu.cn/",
		"http://mirrors.neusoft.edu.cn/",
		"http://mirrors.hit.edu.cn/#/home",
		"https://mirrors.nju.edu.cn/",
		"http://mirrors.cqu.edu.cn/",
		"https://mirrors.dgut.edu.cn/",
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
	errs := LatencyTestWithChan(urls, c, 5, true, 5*time.Second)
	if errs != nil {

	}
	wg.Wait()
	fmt.Println("---------")
}

func TestLatencyTest(t *testing.T) {
	urls := []string{
		"https://mirrors.huaweicloud.com/",
		"https://mirrors.cloud.tencent.com/",
		"http://mirrors.163.com/",
		"http://mirrors.sohu.com/",
		"http://mirrors.yun-idc.com/",
		"https://mirrors.tuna.tsinghua.edu.cn/",
		"http://mirrors.ustc.edu.cn/",
		"http://mirror.neu.edu.cn/",
		"https://mirror.bjtu.edu.cn/",
		"http://mirrors.zju.edu.cn/",
		"http://mirrors.hust.edu.cn/",
		"https://mirrors.njupt.edu.cn/",
		"http://mirrors.neusoft.edu.cn/",
		"http://mirrors.hit.edu.cn/#/home",
		"https://mirrors.nju.edu.cn/",
		"http://mirrors.cqu.edu.cn/",
		"https://mirrors.dgut.edu.cn/",
	}
	results, err := LatencyTest(urls, 5, true, 5*time.Second)
	results.SortByDuration()
	if err != nil {
		log.Println(err)
		return
	}
	for _, v := range results {
		if v.Error != nil {
			color.Red("%s", v.Error)
			continue
		}
		color.Green("%s %d ms", v.Address, v.Latency.Milliseconds())
	}
}
