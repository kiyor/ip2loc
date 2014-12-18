#	ip2loc

##	How to use

```bash

traceroute -n google.com | ip2loc

mtr --no-dns -r -c 10 google.com | ip2loc

```

##	Sample result

![imt](http://ccnacdn.s3.amazonaws.com/img/2014-12-15_nginxln2html__ssh__14144_11-31-44.png)

##	Note

-	change `ip2loc.go` api url to your own [freegeoip](https://github.com/fiorix/freegeoip) server if you want

-	feel free to use or fork