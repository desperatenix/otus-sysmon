package main

import (
	"flag"
	"fmt"
	"github.com/desperatenix/otus-sysmon/api"
	"github.com/desperatenix/otus-sysmon/internal/config"
	"github.com/desperatenix/otus-sysmon/internal/repos"
	"github.com/desperatenix/otus-sysmon/internal/server"
	"github.com/desperatenix/otus-sysmon/internal/sysmon"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	flag.Parse()

	if isVersionCommand() {
		printVersion()
		os.Exit(0)
	}

	mData := &repos.MetricsData{
		Cpu: nil,
		La: nil,
	}

	conf, _ := config.LoadCfg()

	sysmon := &sysmon2.Sysmon2{
		Conf:     &conf,
		MetricsData: mData,
		StopChan: make(chan bool),
		Snapshots: make(map[time.Time]*repos.MetricsData),
	}

	sysmon.Start()

	sysmonServer := &server.SysmonServer{
		Conf:     &conf,
		Stat: sysmon.Snapshots,
		//LARepo:   laRepo,
		//CPURepo:  cpuRepo,
	}

	//sysmonServer.Start

	srv := grpc.NewServer()

	// todo запуск grps server
	port := string(sysmonServer.Conf.Server.Port)
	lsn, err := net.Listen("tcp", ":" + port)
	if err != nil {
		log.Fatalln(err)
	}
	api.RegisterSysMonServer(srv, sysmonServer)
	srv.Serve(lsn)

	// Обработка прерываний
	errListener := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSTOP)
		notifySignal := <-c
		errListener <- fmt.Errorf("%s", notifySignal)
	}()

	err = <-errListener
	srv.Stop()
	sysmon.Stop()
	log.Printf("stop by %s", err)
}
