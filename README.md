#	ip2loc

##	how to use

```bash

traceroute -n google.com | ip2loc

mtr --no-dns -r -c 10 google.com | ip2loc

```

-	change `ip2loc.go` api url to your own [freegeoip](https://github.com/fiorix/freegeoip) server if you want

-	feel free to use or fork
