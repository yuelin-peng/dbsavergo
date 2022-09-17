package dao

import "fmt"

var (
	InvalidParam = fmt.Errorf("invalid order")
	InvalidConn  = fmt.Errorf("invalid connection")
	DBAbnormal   = fmt.Errorf("database abnormal")
	CasNotBeNil  = fmt.Errorf("CAS can't be nil")
)
