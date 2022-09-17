package client

import "context"

type TimeInterval struct {
	StartTime int64
	EndTime   int64 // excluded
}

//go:generate mockgen -destination=mock_db_monitor.go -source=db_monitor_interface.go -package=client DBMonitor

type DBMonitor interface {
	GetDBAbnormalDuration(ctx context.Context) ([]TimeInterval, error)
}
