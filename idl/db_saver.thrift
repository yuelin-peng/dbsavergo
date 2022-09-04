namespace go db_saver_service

struct SaveOrderRequest {
	1:string OrderNO	// 订单号
	2:i64 Version		// 版本号
	3:i64 ModifyTime 	// 修改时间
}

struct SaveOrderResponse {
	1:string OrderNO	// 订单号
	2:i64 Version		// 版本号
	3:string RetCode	// 返回码 0000--代表成功，其他代码失败
	4:string RetMsg		// 返回消息
}

struct QueryOrderRequest {
	1:string OrderNO	// 订单号
}

struct QueryOrderResponse {
	1:string OrderNO	// 订单号
	2:i64 Version		// 版本号
	3:string RetCode	// 返回码 0000--代表成功，其他代码失败
	4:string RetMsg		// 返回消息
}

service DBSaverService {
	// 保存订单
	SaveOrderResponse SaveOrder(1:SaveOrderRequest req)

	// 查询订单
	QueryOrderResponse QueryOrder(1:QueryOrderRequest req)
}

