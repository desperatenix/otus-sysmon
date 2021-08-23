package server

import (
	"fmt"
	"time"

	"github.com/desperatenix/otus-sysmon/api"
	"github.com/desperatenix/otus-sysmon/internal/config"
	"github.com/desperatenix/otus-sysmon/internal/repos"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SysmonServer struct {
	api.UnimplementedSysMonServer
	Conf *config.Config
	Stat repos.Snapshots
}

func (s *SysmonServer) GetStats(data *api.StatsRequest, stream api.SysMon_GetStatsServer) error {
	fmt.Printf("%+v", data)

	n := int(data.N)
	m := int(data.M)

	MaxTimeStorage := s.Conf.Collectors.MaxTimeStorage

	if n <= 0 {
		return status.Error(codes.InvalidArgument, "N must be greater than 0 seconds")
	}
	if n > MaxTimeStorage {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("N must be less than %v seconds", MaxTimeStorage))
	}
	if m <= 0 {
		return status.Error(codes.InvalidArgument, "M must be greater than 0 seconds")
	}
	if m > MaxTimeStorage {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("M must be less than %v seconds", MaxTimeStorage))
	}

	for {
		for t, m := range s.Stat {
			err := stream.Send(&api.Stats{
				Time: timestamppb.New(t),
				LoadAvg: &api.LoadAvg{
					Load1:  m.La.Load1,
					Load5:  m.La.Load5,
					Load15: m.La.Load15,
				},
				Cpu: &api.CPU{
					User:   m.Cpu.User,
					System: m.Cpu.System,
					Idle:   m.Cpu.Idle,
				},
			})
			if err != nil {
				fmt.Println(err)
				return err
			}
		}

		time.Sleep(time.Duration(data.N) * time.Second)
	}
}
