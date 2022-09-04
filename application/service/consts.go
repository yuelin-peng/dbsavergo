package service

type RetInfo struct {
	RetCode string
	RetMsg  string
}

func NewRetInfo(retCode, retMsg string) *RetInfo {
	return &RetInfo{
		RetCode: retCode,
		RetMsg:  retMsg,
	}
}

var (
	SUCCESS      = NewRetInfo("0000", "ok")
	InvalidParam = NewRetInfo("DS0001", "invalid param")
	DBAbnormal   = NewRetInfo("DS0002", "database abnormal")
)
