/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : ip2loc.go

* Purpose :

* Creation Date : 12-14-2014

* Last Modified : Mon 15 Dec 2014 07:24:02 PM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	ipLocMap map[string]*ipLoc
)

type ipLoc struct {
	CountryCode string `json:"country_code"`
	CountryName string `json:"country_name"`
	RegionName  string `json:"region_name"`
	City        string `json:"city"`
}

func ip2loc(ip string) *ipLoc {
	if val, ok := ipLocMap[ip]; ok {
		return val
	}
	var i ipLoc

	res, err := http.Get(fmt.Sprintf("http://66.175.223.83:8080/json/%s", ip))
	if err != nil {
		log.Printf("error %s\n", err.Error())
	}
	b, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(b, &i)

	return &i
}
