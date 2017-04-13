/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : ip2loc.go

* Purpose :

* Creation Date : 12-14-2014

* Last Modified : Thu 13 Apr 2017 04:46:16 PM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	ipLocMap = make(map[string]*ipLoc)
	mu       = new(sync.RWMutex)
	client   = &http.Client{
		Timeout: 3 * time.Second,
	}
)

type ipLoc struct {
	CountryCode string  `json:"country_code"`
	CountryName string  `json:"country_name"`
	RegionName  string  `json:"region_name"`
	City        string  `json:"city"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

func ip2loc(ip string) *ipLoc {
	mu.RLock()
	val, ok := ipLocMap[ip]
	mu.RUnlock()

	id := random(4)

	if ok {
		if *enableMap {
			mapData.Lock()
			mapData.m[mapKey{time.Now(), id}] = val
			mapData.Unlock()
		}
		return val
	}

	var i ipLoc

	req, _ := http.NewRequest("GET", fmt.Sprintf("http://%s/json/%s", *flagIp, ip), nil)
	req.Host = *flagHost
	res, err := client.Do(req)
	if err != nil {
		log.Printf("error %s\n", err.Error())
		return nil
	}
	b, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(b, &i)

	mu.Lock()
	defer mu.Unlock()
	ipLocMap[ip] = &i

	if *enableMap {
		mapData.Lock()
		mapData.m[mapKey{time.Now(), id}] = &i
		mapData.Unlock()
	}

	return &i
}
