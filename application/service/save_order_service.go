package service

import (
	"context"
	"fmt"

	"github.com/yuelin-peng/dbsavergo/kitex_gen/db_saver_service"
)

func SaveOrder(ctx context.Context, req *db_saver_service.SaveOrderRequest) (*db_saver_service.SaveOrderResponse, error) {
	if err := checkSaveOrderParam(req); err != nil {
		return &db_saver_service.SaveOrderResponse{
			RetCode: InvalidParam.RetCode,
			RetMsg:  InvalidParam.RetMsg,
		}, nil
	}
	orderModel, err := createOrderModelForSave(ctx)
	if err != nil || orderModel == nil {
		return &db_saver_service.SaveOrderResponse{
			OrderNO: req.OrderNO,
			Version: req.Version,
			RetCode: DBAbnormal.RetCode,
			RetMsg:  DBAbnormal.RetMsg,
		}, nil
	}
	return &db_saver_service.SaveOrderResponse{
		OrderNO: req.OrderNO,
		Version: req.Version,
		RetCode: SUCCESS.RetCode,
		RetMsg:  SUCCESS.RetMsg,
	}, nil
}

func checkSaveOrderParam(req *db_saver_service.SaveOrderRequest) error {
	if req == nil {
		return fmt.Errorf("request can't be nil")
	}
	if len(req.OrderNO) == 0 {
		return fmt.Errorf("order number can't be empty")
	}
	if req.Version < 0 {
		return fmt.Errorf("version can't be less than 0")
	}

	return nil
}
