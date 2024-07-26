package main

import (
	"flag"
	"fmt"
	"github.com/xueqianLu/vehackcenter/config"
	"github.com/xueqianLu/vehackcenter/server"
)

var (
	servicePort = flag.Int("port", 9000, "service port")
	hackerCount = flag.Int("c", 33, "hacker count")
	vote        = flag.Int("vote", 1, "vote value, 1 or 0")
	beginHeight = flag.Int("begin", 500, "the height begin to execute hack")
)

func main() {
	flag.Parse()
	url := fmt.Sprintf("0.0.0.0:%d", *servicePort)
	if *vote != 0 {
		*vote = 1
	}
	conf := config.Config{
		Url:         url,
		HackerCount: *hackerCount,
		BeginToHack: *beginHeight,
		Vote:        *vote,
	}

	node := server.NewNode(conf)
	node.RunServer()
}
