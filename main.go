package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/xueqianLu/vehackcenter/config"
	"github.com/xueqianLu/vehackcenter/server"
	"io"
	"os"
)

var (
	servicePort = flag.Int("port", 9000, "service port")
	hackerCount = flag.Int("c", 33, "hacker count")
	vote        = flag.Int("vote", 1, "vote value, 1 or 0")
	beginHeight = flag.Int("begin", 500, "the height begin to execute hack")
	endHeight   = flag.Int("end", 100000, "the height end to execute hack")
	logfile     = flag.String("log", "", "log file path")
)

func main() {
	flag.Parse()
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, TimestampFormat: "2006-01-02 15:04:05.000"})
	var logfileCtx *os.File
	var err error
	if *logfile != "" {
		logfileCtx, err = os.OpenFile(*logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			log.SetOutput(io.MultiWriter(logfileCtx, os.Stdout))
		} else {
			log.Info("Failed to log to file, using default stderr")
		}
	}
	defer func() {
		if logfileCtx != nil {
			logfileCtx.Close()
		}
	}()

	url := fmt.Sprintf("0.0.0.0:%d", *servicePort)
	if *vote != 0 {
		*vote = 1
	}
	conf := config.Config{
		Url:         url,
		HackerCount: *hackerCount,
		BeginToHack: *beginHeight,
		EndToHack:   *endHeight,
		Vote:        *vote,
	}

	node := server.NewNode(conf)
	node.RunServer()
}
