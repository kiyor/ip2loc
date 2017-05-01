#	ip2loc

1. If you haven't install golang then [download and install](http://golang.org/dl/)
2. `go get -u github.com/kiyor/ip2loc`
3. `go install github.com/kiyor/ip2loc`

##	How to use

```bash

traceroute -n google.com | ip2loc

mtr --no-dns -r -c 10 google.com | ip2loc

```

##	Sample result

![imt](http://ccnacdn.s3.amazonaws.com/img/2014-12-15_nginxln2html__ssh__14144_11-31-44.png)

### Enable map output

run:

```bash
sudo cat /var/log/secure|ip2loc -map >/dev/null
```

open browser http://${ip}:7676

![](https://kiyor.s3.amazonaws.com/imgs/2017-05-01_15-19-36_Marker_Clustering.png)

##	Note

-	change `ip2loc.go` api url to your own [freegeoip](https://github.com/fiorix/freegeoip) server if you want

-	feel free to use or fork
