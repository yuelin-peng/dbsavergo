package main

import (
	"context"
	"workspace/dbsavergo/kitex_gen/db_saver_service"
)

// DBSaverServiceImpl implements the last service interface defined in the IDL.
type DBSaverServiceImpl struct{}

// SaveOrder implements the DBSaverServiceImpl interface.
func (s *DBSaverServiceImpl) SaveOrder(ctx context.Context, req *db_saver_service.SaveOrderRequest) (resp *db_saver_service.SaveOrderResponse, err error) {
	// TODO: Your code here...
	return
}

// QueryOrder implements the DBSaverServiceImpl interface.
func (s *DBSaverServiceImpl) QueryOrder(ctx context.Context, req *db_saver_service.QueryOrderRequest) (resp *db_saver_service.QueryOrderResponse, err error) {
	// TODO: Your code here...
	return
}
