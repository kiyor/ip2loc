/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : main.go

* Purpose :

* Creation Date : 12-14-2014

* Last Modified : Thu 13 Apr 2017 02:21:46 AM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/wsxiaoys/terminal/color"
	"io"
	"log"
	"net"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	reIp             = regexp.MustCompile(`(\d+\.\d+\.\d+\.\d+)`)
	enableLL   *bool = flag.Bool("ll", false, "enable show longitude and latitude")
	enableMap  *bool = flag.Bool("map", false, "enable google map")
	flagIp           = flag.String("addr", "ip.2ns.io", "freegeoip address")
	flagHost         = flag.String("host", *flagIp, "freegeoip Host header address")
	flagExpire       = flag.Duration("expire", 15*time.Minute, "expire time")
	// 	cache     gcache.Cache
	mapData = newMapData()
)

type MapData struct {
	m map[mapKey]*ipLoc
	*sync.RWMutex
}

type mapKey struct {
	t time.Time
	r string
}

func newMapData() MapData {
	return MapData{
		m:       make(map[mapKey]*ipLoc),
		RWMutex: new(sync.RWMutex),
	}
}

func cron() {
	ticker := time.NewTicker(time.Minute * 1)
	go func() {
		for range ticker.C {
			cleanMap()
		}
	}()
}

func cleanMap() {
	mapData.Lock()
	defer mapData.Unlock()
	for k := range mapData.m {
		if time.Now().Sub(k.t) > *flagExpire {
			delete(mapData.m, k)
		}
	}
}

func init() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
}

type Line struct {
	index int64
	line  string
}

type Wg struct {
	sync.WaitGroup
	count *int64
}

func (wg *Wg) add(i int) {
	wg.Add(i)
	*wg.count++
}

func (wg *Wg) done() {
	*wg.count--
	wg.Done()
}

func main() {
	if *enableMap {
		go runHttp()
	}

	stop := make(chan bool)
	ch := make(chan Line, 1)
	var wg Wg
	wg.count = new(int64)

	var i int64
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			l, err := reader.ReadString('\n')

			if err != nil {
				if err == io.EOF {
					wg.Wait()
					stop <- true
				} else {
					log.Println(err.Error())
					os.Exit(1)
				}
			} else {
				line := Line{
					index: i,
					line:  l,
				}
				for *wg.count > int64(runtime.NumCPU()) {
					time.Sleep(1 * time.Millisecond)
				}
				wg.add(1)
				go processing(line, ch)
				i++
			}
		}
	}()

	var i2 int64
	for {
		select {
		case l := <-ch:
			switch l.index {
			case i2:
				fmt.Println(l.line)
				i2++
				wg.done()
			default:
				// if index is not expect, then create a backgroup process send back channel
				go func(l Line) {
					time.Sleep(1 * time.Millisecond)
					ch <- l
				}(l)
			}
		case <-stop:
			os.Exit(0)
		}
	}
}

func processing(line Line, ch chan Line) {
	if len(line.line) > 0 {
		line.line = strings.Trim(line.line, "\n")
	}
	if reIp.MatchString(line.line) {
		part := reIp.FindAllStringSubmatch(line.line, -1)

		done := make(map[string]bool)
		for _, v := range part {
			ipStr := v[1]
			if ip := net.ParseIP(ipStr); ip != nil {
				if ip.IsLoopback() {
					continue
				}
				if _, ok := done[ipStr]; !ok {
					loc := ip2loc(ipStr)
					var replace string
					if len(loc.CountryCode) > 0 {
						replace += loc.CountryCode + " "
					}
					if len(loc.RegionName) > 0 {
						replace += loc.RegionName + " "
					}
					if len(loc.City) > 0 {
						replace += loc.City + " "
					}
					if *enableLL {
						if loc.Latitude != 0 {
							replace += fmt.Sprintf("%v %v ", loc.Latitude, loc.Longitude)
						}
					}
					if len(replace) > 0 {
						replace = replace[:len(replace)-1]
						line.line = strings.Replace(line.line, ipStr, color.Sprintf("@{m}%v@{|}@{r}[@{g}%v@{|}@{r}]@{|}", ip, replace), -1)
					}
					done[ipStr] = true
				}
			}
		}

		ch <- line
		return
	} else {
		ch <- line
		return
	}
}
