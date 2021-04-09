package filter

import (
	"errors"
	"net"
	"strings"
	"sync"
	"time"

	colorful "github.com/fatih/color"
)

// SpeedTest test latency between mirror site and local
// store valid mirror's index in (chan int)
func SpeedTest(urls []string, c chan<- int) error {
	if len(urls) < 1 {
		return errors.New("mirror sites must more than one")
	}
	fixURLs := urlHandler(urls)
	var wg sync.WaitGroup
	for idx, url := range fixURLs {
		wg.Add(1)
		go func(url string, idx int) {
			defer wg.Done()
			_, err := net.DialTimeout("tcp", url+":http", 5*time.Second)
			if err != nil {
				colorful.Red("skiped mirror: %s\ndue to %s\n", url, err.Error())
				return
			}
			c <- idx
		}(url, idx)
	}
	wg.Wait()
	return nil
}

// to get domain of url
func urlHandler(urls []string) (validURLs []string) {
	for _, url := range urls {
		s := strings.Replace(url, "https://", "", 1)
		s = strings.Replace(s, "http://", "", 1)
		validURLs = append(validURLs, s)
	}
	return
}
