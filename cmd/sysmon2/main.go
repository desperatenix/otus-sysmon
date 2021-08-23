package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/desperatenix/otus-sysmon/api"
	"github.com/desperatenix/otus-sysmon/internal/config"
	"github.com/desperatenix/otus-sysmon/internal/repos"
	"github.com/desperatenix/otus-sysmon/internal/server"
	"github.com/desperatenix/otus-sysmon/internal/sysmon"
	"google.golang.org/grpc"
)

var cfgPath string

func init() {
	flag.StringVar(&cfgPath, "configPath", "./cfg/config.yml", "Path to configuration file without name.")
}

func main() {
	flag.Parse()

	if isVersionCommand() {
		printVersion()
		os.Exit(0)
	}

	mData := &repos.MetricsData{
		Cpu: nil,
		La:  nil,
	}

	conf, _ := config.LoadCfg(cfgPath)

	sysmon := &sysmon2.Sysmon2{
		Conf:        &conf,
		MetricsData: mData,
		StopChan:    make(chan bool),
		Snapshots:   make(map[time.Time]*repos.MetricsData),
	}

	sysmon.Start()

	sysmonServer := &server.SysmonServer{
		Conf: &conf,
		Stat: sysmon.Snapshots,
	}

	srv := grpc.NewServer()

	port := sysmonServer.Conf.Server.Port
	lsn, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalln(err)
	}
	api.RegisterSysMonServer(srv, sysmonServer)
	err = srv.Serve(lsn)
	if err != nil {
		log.Fatalln(err)
	}

	// Обработка прерываний
	errListener := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSTOP) //nolint:staticcheck
		notifySignal := <-c
		errListener <- fmt.Errorf("%s", notifySignal)
	}()

	err = <-errListener
	srv.Stop()
	sysmon.Stop()
	log.Printf("stop by %s", err)
}
